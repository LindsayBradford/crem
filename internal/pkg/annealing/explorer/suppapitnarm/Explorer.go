// Copyright (c) 2018 Australian Rivers Institute.

package suppapitnarm

import (
	"fmt"
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

	observer.SynchronousAnnealingEventNotifier
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
	ke.note("Trying Random Model Change")

	compressedInitialModelState := ke.modelArchive.Compress(ke.Model())
	ke.Model().DoRandomChange()
	compressedChangedModelState := ke.modelArchive.Compress(ke.Model())

	variableDifferences := compressedChangedModelState.VariableDifferences(compressedInitialModelState)

	event := observer.NewEvent(observer.Explorer).
		WithNote("Attempting to Archive Changed Model").
		WithAttribute("Model SHA256", compressedChangedModelState.Sha256())

	ke.NotifyObserversOfEvent(*event)

	ke.archiveStorageResult = ke.modelArchive.AttemptToArchiveState(compressedChangedModelState)

	ke.AcceptOrRevertChange(variableDifferences)
	ke.ReturnToBaseIfRequired(compressedChangedModelState)
}

func (ke *Explorer) note(note string) {
	noteEvent := observer.NewEvent(observer.Explorer).
		WithAttribute(observer.Note.String(), note)
	ke.NotifyObserversOfEvent(*noteEvent)
}

func (ke *Explorer) AcceptOrRevertChange(variableDifferences []float64) {
	if ke.changeTriedIsDesirable() {
		ke.setAcceptanceProbability(explorer.Guaranteed)
		ke.changeAccepted = true
		ke.notifyDesirableAcceptance()
	} else {
		if ke.coolant.DecideIfAcceptable(variableDifferences) {
			ke.notifyUndesirableAcceptance()
			ke.AcceptLastChange()
		} else {
			ke.notifyUndesirableReversion()
			ke.RevertLastChange()
		}
	}
}

func (ke *Explorer) notifyDesirableAcceptance() {
	event := observer.NewEvent(observer.Explorer).
		WithNote("Accepting Desirable Change").
		WithAttribute(explorer.AcceptanceProbability, ke.coolant.AcceptanceProbability())

	ke.NotifyObserversOfEvent(*event)
}

func (ke *Explorer) notifyUndesirableAcceptance() {
	acceptEvent := observer.NewEvent(observer.Explorer).
		WithNote("Accepting Undesirable Change").
		WithAttribute(explorer.AcceptanceProbability, ke.coolant.AcceptanceProbability())

	ke.NotifyObserversOfEvent(*acceptEvent)
}

func (ke *Explorer) notifyUndesirableReversion() {
	revertEvent := observer.NewEvent(observer.Explorer).
		WithNote("Reverting Undesirable Change").
		WithAttribute(explorer.AcceptanceProbability, ke.coolant.AcceptanceProbability())

	ke.NotifyObserversOfEvent(*revertEvent)
}

func (ke *Explorer) changeTriedIsDesirable() bool {
	switch ke.archiveStorageResult {
	case archive.StoredWithNoDominanceDetected, archive.StoredReplacingDominatedEntries,
		archive.RejectedWithDuplicateEntryDetected:
		ke.changeIsDesirable = true
	case archive.RejectedWithStoredEntryDominanceDetected:
		ke.changeIsDesirable = false
	}

	revertEvent := observer.NewEvent(observer.Explorer).
		WithAttribute("ArchiveStorageResult", ke.archiveStorageResult.String()).
		WithAttribute("ChangeDesirable", ke.changeIsDesirable)

	ke.NotifyObserversOfEvent(*revertEvent)

	return ke.changeIsDesirable
}

func (ke *Explorer) AcceptLastChange() {
	ke.archiveStorageResult = ke.modelArchive.ForceIntoArchive(ke.Model())
	ke.changeAccepted = true

	event := observer.NewEvent(observer.Explorer).
		WithNote("Forcing Model into Archive").
		WithAttribute("ArchiveStorageResult", ke.archiveStorageResult.String())

	ke.NotifyObserversOfEvent(*event)
}

func (ke *Explorer) RevertLastChange() {
	ke.Model().UndoChange()
	ke.changeAccepted = false
}

func (ke *Explorer) ReturnToBaseIfRequired(state *archive.CompressedModelState) {
	if !ke.shouldReturnToBase() {
		return
	}
	ke.returnToBase(state)
	ke.adjustReturnToBaseRate()
}

func (ke *Explorer) shouldReturnToBase() bool {
	ke.iterationsUntilReturnToBase--
	return ke.iterationsUntilReturnToBase <= 0
}

func (ke *Explorer) returnToBase(currentModelState *archive.CompressedModelState) {
	const minFraction = 1e-63
	selectionRangeLimit := int(math.Ceil(float64(ke.modelArchive.Len()) * ke.returnToBaseIsolationFraction))

	compressedModel := ke.modelArchive.SelectRandomIsolatedModel(selectionRangeLimit)

	// Is this actually a problem?  Suppapitnarm paper doesn't mention it.  Original CRP doesn't cater for it.
	// It just seems odd to trawl through the isolated entries only to return to the current if its isolated.
	if currentModelState.Sha256() == compressedModel.Sha256() {
		warningMessage := fmt.Sprintf("Randomly selected return-to-base isolated model is same as current [%s]",
			currentModelState.Sha256())
		ke.LogHandler().Warn(warningMessage)
	}

	ke.modelArchive.Decompress(compressedModel, ke.Model())

	if ke.returnToBaseIsolationFraction == 1 {
		ke.returnToBaseIsolationFraction = ke.parameters.GetFloat64(ReturnToBaseIsolationFraction)
	} else {
		ke.returnToBaseIsolationFraction *= ke.returnToBaseIsolationFraction
		ke.returnToBaseIsolationFraction = math.Max(ke.returnToBaseIsolationFraction, minFraction)
	}

	event := observer.NewEvent(observer.Explorer).
		WithNote("Returning to Base").
		WithAttribute("SelectionRangeLimit", selectionRangeLimit).
		WithAttribute("IsolationFraction", ke.returnToBaseIsolationFraction).
		WithAttribute("New Base Model SHA256", compressedModel.Sha256())

	ke.NotifyObserversOfEvent(*event)
}

func (ke *Explorer) adjustReturnToBaseRate() {
	minimumBaseRate := float64(ke.parameters.GetInt64(MinimumReturnToBaseRate))
	rateAdjustment := ke.parameters.GetFloat64(ReturnToBaseAdjustmentFactor)

	ke.returnToBaseStep = math.Max(minimumBaseRate, ke.returnToBaseStep*rateAdjustment)
	ke.deriveIterationsUntilReturnToBase()
}

func (ke *Explorer) deriveIterationsUntilReturnToBase() {
	ke.iterationsUntilReturnToBase = uint64(ke.returnToBaseStep)

	event := observer.NewEvent(observer.Explorer).
		WithAttribute("IterationsUntilReturnToBase", ke.iterationsUntilReturnToBase)
	ke.NotifyObserversOfEvent(*event)
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
			Add(ArchiveSize, ke.modelArchive.Len())
	}
	return nil
}

func (ke *Explorer) CoolDown() {
	ke.coolant.CoolDown()
	ke.notifyCoolDown()
}

func (ke *Explorer) notifyCoolDown() {
	event := observer.NewEvent(observer.Explorer).
		WithNote("Cooling").
		WithAttribute(explorer.Temperature, ke.coolant.Temperature())

	ke.NotifyObserversOfEvent(*event)
}
