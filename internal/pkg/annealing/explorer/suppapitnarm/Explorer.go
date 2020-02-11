// Copyright (c) 2018 Australian Rivers Institute.

package suppapitnarm

import (
	"math"
	"strings"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/cooling"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/cooling/coolants/suppapitnarm"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/archive"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/attributes"
	errors2 "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/pkg/errors"
)

const (
	nameSeparator = ","

	ArchiveSize                     = "ArchiveSize"
	ArchiveResult                   = "ArchiveResult"
	IterationsUntilNextReturnToBase = "IterationsUntilNextReturnToBase"
)

type Explorer struct {
	name.NameContainer
	name.IdentifiableContainer

	model.ContainedModel
	loggers.ContainedLogger

	coolant cooling.TemperatureCoolant

	scenarioId string

	parameters            Parameters
	optimisationDirection optimisationDirection
	objectiveVariableName string

	modelArchive         archive.NonDominanceModelArchive
	archiveStorageResult archive.StorageResult

	iterationsUntilReturnToBase   uint64
	returnToBaseStep              float64
	returnToBaseIsolationFraction float64

	changeIsDesirable    bool
	changeAccepted       bool
	objectiveValueChange float64
}

func New() *Explorer {
	newExplorer := new(Explorer)
	newExplorer.parameters.Initialise()
	newExplorer.SetCoolant(suppapitnarm.NewCoolant())
	newExplorer.modelArchive.Initialise()
	newExplorer.SetModel(model.NewNullModel())
	return newExplorer
}

func (ke *Explorer) Initialise() {
	ke.LogHandler().Debug(ke.scenarioId + ": Initialising Solution Explorer")
	ke.modelArchive.Initialise()
	ke.coolant.SetRandomNumberGenerator(rand.NewTimeSeeded())
	ke.Model().Initialise()
	ke.deriveIterationsUntilReturnToBase()
}

func (ke *Explorer) WithName(name string) *Explorer {
	ke.SetName(name)
	return ke
}

func (ke *Explorer) WithModel(model model.Model) *Explorer {
	ke.SetModel(model)
	return ke
}

func (ke *Explorer) WithCoolant(coolant cooling.TemperatureCoolant) *Explorer {
	ke.SetCoolant(coolant)
	return ke
}

func (ke *Explorer) SetCoolant(coolant cooling.TemperatureCoolant) {
	ke.coolant = coolant
}

func (ke *Explorer) WithParameters(params parameters.Map) *Explorer {
	ke.SetParameters(params)
	return ke
}

func (ke *Explorer) SetParameters(params parameters.Map) error {
	ke.parameters.AssignOnlyEnforcedUserValues(params)
	ke.coolant.SetParameters(params)

	ke.returnToBaseStep = float64(ke.parameters.GetInt64(InitialReturnToBaseStep))
	ke.returnToBaseIsolationFraction = 1

	ke.checkDecisionVariablesFromParams()

	return ke.parameters.ValidationErrors()
}

func (ke *Explorer) SetTemperature(temperature float64) error {
	if temperature <= 0 {
		return errors.New("invalid attempt to set annealer temperature to value <= 0")
	}
	ke.coolant.SetTemperature(temperature)
	return nil
}

func (ke *Explorer) checkDecisionVariablesFromParams() {
	decisionVariableNames := ke.parameters.GetString(ExplorableDecisionVariables)
	splitVariableNames := strings.Split(decisionVariableNames, nameSeparator)
	for _, name := range splitVariableNames {
		variableOffered := ke.Model().OffersDecisionVariable(name)
		if !variableOffered {
			ke.parameters.AddValidationErrorMessage("decision variable [" + name + "] not recognised by model")
		}
	}
}

func (ke *Explorer) ParameterErrors() error {
	mergedErrors := errors2.New("Kirkpatrick Explorer Parameter Validation")

	mergedErrors.Add(ke.parameters.ValidationErrors())
	mergedErrors.Add(ke.coolant.ParameterErrors())

	if mergedErrors.Size() > 0 {
		return mergedErrors
	}

	return nil
}

func (ke *Explorer) ObjectiveValue() float64 {
	variable := ke.Model().DecisionVariable(ke.objectiveVariableName)
	return variable.Value()
}

func (ke *Explorer) TryRandomChange() {
	compressedInitialModelState := ke.modelArchive.Compress(ke.Model())
	ke.Model().DoRandomChange()
	compressedChangedModelState := ke.modelArchive.Compress(ke.Model())
	variableDifferences := compressedChangedModelState.VariableDifferences(compressedInitialModelState)
	ke.archiveStorageResult = ke.modelArchive.AttemptToArchiveState(compressedChangedModelState)
	ke.AcceptOrRevertChange(variableDifferences)
	ke.ReturnToBaseIfRequired()
}

func (ke *Explorer) AcceptOrRevertChange(variableDifferences []float64) {
	if ke.changeTriedIsDesirable() {
		ke.setAcceptanceProbability(explorer.Guaranteed)
		ke.changeAccepted = true
	} else {
		if ke.coolant.DecideIfAcceptable(variableDifferences) {
			ke.AcceptLastChange()
		} else {
			ke.RevertLastChange()
		}
	}
}

func (ke *Explorer) changeTriedIsDesirable() bool {
	switch ke.archiveStorageResult {
	case archive.StoredWithNoDominanceDetected, archive.StoredReplacingDominatedEntries:
		ke.changeIsDesirable = true
	case archive.RejectedWithStoredEntryDominanceDetected, archive.RejectedWithDuplicateEntryDetected:
		ke.changeIsDesirable = false
	}
	return ke.changeIsDesirable
}

func (ke *Explorer) AcceptLastChange() {
	ke.archiveStorageResult = ke.modelArchive.ForceIntoArchive(ke.Model())
	ke.changeAccepted = true
}

func (ke *Explorer) RevertLastChange() {
	ke.Model().UndoChange()
	ke.changeAccepted = false
}

func (ke *Explorer) ReturnToBaseIfRequired() {
	if !ke.shouldReturnToBase() {
		return
	}
	ke.returnToBase()
	ke.adjustReturnToBaseRate()
}

func (ke *Explorer) shouldReturnToBase() bool {
	ke.iterationsUntilReturnToBase--
	return ke.iterationsUntilReturnToBase <= 0
}

func (ke *Explorer) returnToBase() {
	const minFraction = 1e-63
	selectionRangeLimit := int(math.Ceil(float64(ke.modelArchive.Len()) * ke.returnToBaseIsolationFraction))
	compressedModel := ke.modelArchive.SelectRandomIsolatedModel(selectionRangeLimit)
	ke.modelArchive.Decompress(compressedModel, ke.Model())
	if ke.returnToBaseIsolationFraction == 1 {
		ke.returnToBaseIsolationFraction = ke.parameters.GetFloat64(ReturnToBaseIsolationFraction)
	} else {
		ke.returnToBaseIsolationFraction *= ke.returnToBaseIsolationFraction
		ke.returnToBaseIsolationFraction = math.Max(ke.returnToBaseIsolationFraction, minFraction)
	}
}

func (ke *Explorer) adjustReturnToBaseRate() {
	minimumBaseRate := float64(ke.parameters.GetInt64(MinimumReturnToBaseRate))
	rateAdjustment := ke.parameters.GetFloat64(ReturnToBaseAdjustmentFactor)

	ke.returnToBaseStep = math.Max(minimumBaseRate, ke.returnToBaseStep*rateAdjustment)
	ke.deriveIterationsUntilReturnToBase()
}

func (ke *Explorer) deriveIterationsUntilReturnToBase() {
	ke.iterationsUntilReturnToBase = uint64(ke.returnToBaseStep)
}

func (ke *Explorer) DeepClone() explorer.Explorer {
	clone := *ke
	clone.coolant.SetRandomNumberGenerator(rand.NewTimeSeeded())
	modelClone := ke.Model().DeepClone()
	clone.SetModel(modelClone)
	return &clone
}

func (ke *Explorer) TearDown() {
	ke.LogHandler().Debug(ke.scenarioId + ": Triggering tear-down of Solution Explorer")
	ke.Model().TearDown()
}

func (ke *Explorer) setAcceptanceProbability(probability float64) {
	ke.coolant.SetAcceptanceProbability(math.Min(explorer.Guaranteed, probability))
}

func (ke *Explorer) EventAttributes(eventType observer.EventType) attributes.Attributes {
	baseAttributes := new(attributes.Attributes).
		Add(explorer.Temperature, ke.coolant.Temperature())
	switch eventType {
	case observer.StartedAnnealing:
		return baseAttributes.Add(explorer.CoolingFactor, ke.coolant.CoolingFactor())
	case observer.StartedIteration, observer.FinishedAnnealing:
		return baseAttributes.Add(ArchiveSize, ke.modelArchive.Len())
	case observer.FinishedIteration:
		return baseAttributes.
			Add(ArchiveSize, ke.modelArchive.Len()).
			Add(ArchiveResult, ke.archiveStorageResult).
			Add(IterationsUntilNextReturnToBase, ke.iterationsUntilReturnToBase).
			Add(explorer.ChangeIsDesirable, ke.changeIsDesirable).
			Add(explorer.AcceptanceProbability, ke.coolant.AcceptanceProbability()).
			Add(explorer.ChangeAccepted, ke.changeAccepted)
	}
	return nil
}

func (ke *Explorer) CoolDown() {
	ke.coolant.CoolDown()
}
