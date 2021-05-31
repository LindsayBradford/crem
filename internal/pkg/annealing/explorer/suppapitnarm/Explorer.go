// Copyright (c) 2018 Australian Rivers Institute.

package suppapitnarm

import (
	"fmt"
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
	"math"
)

const (
	nameSeparator = ","

	ArchiveSize                     = "ArchiveSize"
	ArchiveResult                   = "ArchiveResult"
	IterationsUntilNextReturnToBase = "IterationsUntilNextReturnToBase"
	ModelArchive                    = "ModelArchive"
	LastReturnedToBase              = "LastReturnedToBase"
)

type Explorer struct {
	name.NameContainer
	name.IdentifiableContainer

	currentModel   model.Model
	potentialModel model.Model

	loggers.ContainedLogger

	coolant cooling.TemperatureCoolant

	scenarioId string

	parameters            Parameters
	optimisationDirection optimisationDirection
	objectiveVariableName string

	modelArchive         archive.NonDominanceModelArchive
	archiveStorageResult archive.StorageResult

	currentIteration   uint64
	lastReturnedToBase uint64

	iterationsUntilReturnToBase   uint64
	returnToBaseStep              float64
	returnToBaseIsolationFraction float64

	changeIsDesirable    bool
	changeAccepted       bool
	objectiveValueChange float64

	observer.SynchronousAnnealingEventNotifier
	baseAttributes attributes.Attributes

	desirableAcceptanceEvent   *observer.Event
	undesirableAcceptanceEvent *observer.Event
	undesirableReversionEvent  *observer.Event
	noteEvent                  *observer.Event
}

func New() *Explorer {
	newExplorer := new(Explorer)
	newExplorer.parameters.Initialise()
	newExplorer.SetCoolant(suppapitnarm.NewCoolant())
	newExplorer.modelArchive.Initialise()

	newExplorer.currentModel = model.NewNullModel()
	newExplorer.potentialModel = model.NewNullModel()

	return newExplorer
}

func (ke *Explorer) Initialise() {
	ke.LogHandler().Debug(ke.scenarioId + ": Initialising Solution Explorer")
	ke.modelArchive.Initialise()
	ke.coolant.SetRandomNumberGenerator(rand.NewTimeSeeded())

	ke.currentModel.Initialise()
	ke.currentModel.Randomize()

	ke.potentialModel.Initialise()

	ke.deriveIterationsUntilReturnToBase()
	ke.currentIteration = 1

	ke.baseAttributes = new(attributes.Attributes).
		Add(explorer.Temperature, ke.coolant.Temperature()).
		Add(ArchiveSize, ke.modelArchive.Len())

	ke.desirableAcceptanceEvent = observer.NewEvent(observer.Explorer).
		WithNote("Accepting Desirable Change").
		WithAttribute(explorer.AcceptanceProbability, ke.coolant.AcceptanceProbability())

	ke.undesirableAcceptanceEvent = observer.NewEvent(observer.Explorer).
		WithNote("Accepting Undesirable Change").
		WithAttribute(explorer.AcceptanceProbability, ke.coolant.AcceptanceProbability())

	ke.undesirableReversionEvent = observer.NewEvent(observer.Explorer).
		WithNote("Reverting Undesirable Change").
		WithAttribute(explorer.AcceptanceProbability, ke.coolant.AcceptanceProbability())

	ke.noteEvent = observer.NewEvent(observer.Explorer).
		WithAttribute(observer.Note.String(), "")
}

func (ke *Explorer) WithName(name string) *Explorer {
	ke.SetName(name)
	return ke
}

func (ke *Explorer) WithModel(model model.Model) *Explorer {
	ke.currentModel = model
	ke.potentialModel = model.DeepClone()
	return ke
}

func (ke *Explorer) WithCoolant(coolant cooling.TemperatureCoolant) *Explorer {
	ke.SetCoolant(coolant)
	return ke
}

func (ke *Explorer) SetCoolant(coolant cooling.TemperatureCoolant) {
	ke.coolant = coolant
}

func (ke *Explorer) SetId(id string) {
	ke.IdentifiableContainer.SetId(id)
	ke.modelArchive.SetId(id)
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

	ke.baseAttributes = new(attributes.Attributes).
		Add(explorer.Temperature, ke.coolant.Temperature()).
		Add(ArchiveSize, ke.modelArchive.Len())

	return ke.parameters.ValidationErrors()
}

func (ke *Explorer) SetTemperature(temperature float64) error {
	if temperature <= 0 {
		return errors.New("invalid attempt to set annealer temperature to value <= 0")
	}
	ke.coolant.SetTemperature(temperature)
	return nil
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
	variable := ke.currentModel.DecisionVariable(ke.objectiveVariableName)
	return variable.Value()
}

func (ke *Explorer) TryRandomChange() {
	ke.note("Trying Random Model Change")

	compressedInitialModelState := ke.modelArchive.Compress(ke.currentModel)
	ke.generatePotentialModel()
	compressedChangedModelState := ke.modelArchive.Compress(ke.potentialModel)

	variableDifferences := compressedChangedModelState.VariableDifferences(compressedInitialModelState)

	event := observer.NewEvent(observer.Explorer).
		WithNote("Attempting to Archive Changed Model").
		WithAttribute("Model Encoding", compressedChangedModelState.Encoding())

	ke.NotifyObserversOfEvent(*event)

	ke.archiveStorageResult = ke.modelArchive.AttemptToArchiveState(compressedChangedModelState)

	ke.AcceptOrRevertChange(variableDifferences)
	ke.ReturnToBaseIfRequired(compressedChangedModelState)

	ke.currentIteration++
}

func (ke *Explorer) generatePotentialModel() {
	ke.note("Creating and Randomizing potential new model off old.")
	ke.potentialModel.SynchroniseTo(ke.currentModel)
	ke.potentialModel.Randomize()
	ke.note("Finished creating and Randomizing potential new model.")
}

func (ke *Explorer) note(note string) {
	ke.noteEvent.ReplaceAttribute(observer.Note.String(), note)
	ke.NotifyObserversOfEvent(*ke.noteEvent)
}

func (ke *Explorer) AcceptOrRevertChange(variableDifferences []float64) {
	if ke.changeTriedIsDesirable() {
		ke.AcceptDesirableChange()
		ke.notifyDesirableAcceptance()
	} else {
		if ke.coolant.DecideIfAcceptable(variableDifferences) {
			ke.notifyUndesirableAcceptance()
			ke.AcceptUndesirableChange()
		} else {
			ke.notifyUndesirableReversion()
			ke.RevertLastChange()
		}
	}
}

func (ke *Explorer) AcceptDesirableChange() {
	ke.setAcceptanceProbability(explorer.Guaranteed)
	ke.changeAccepted = true
	ke.currentModel.SynchroniseTo(ke.potentialModel)
}

func (ke *Explorer) notifyDesirableAcceptance() {
	ke.desirableAcceptanceEvent.ReplaceAttribute(explorer.AcceptanceProbability, ke.coolant.AcceptanceProbability())
	ke.NotifyObserversOfEvent(*ke.desirableAcceptanceEvent)
}

func (ke *Explorer) notifyUndesirableAcceptance() {
	ke.undesirableAcceptanceEvent.ReplaceAttribute(explorer.AcceptanceProbability, ke.coolant.AcceptanceProbability())
	ke.NotifyObserversOfEvent(*ke.undesirableAcceptanceEvent)
}

func (ke *Explorer) notifyUndesirableReversion() {
	ke.undesirableReversionEvent.ReplaceAttribute(explorer.AcceptanceProbability, ke.coolant.AcceptanceProbability())
	ke.NotifyObserversOfEvent(*ke.undesirableReversionEvent)
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

func (ke *Explorer) AcceptUndesirableChange() {
	ke.archiveStorageResult = ke.modelArchive.ForceIntoArchive(ke.potentialModel)
	ke.changeAccepted = true
	ke.currentModel.SynchroniseTo(ke.potentialModel)

	event := observer.NewEvent(observer.Explorer).
		WithNote("Forcing Model into Archive").
		WithAttribute("ArchiveStorageResult", ke.archiveStorageResult.String())

	ke.NotifyObserversOfEvent(*event)
}

func (ke *Explorer) RevertLastChange() {
	// we just ignore potential model state.
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
	ke.returnToBaseInsideArchive(currentModelState)
	//ke.returnToBaseOutsideArchive(currentModelState)
}

func (ke *Explorer) returnToBaseInsideArchive(currentModelState *archive.CompressedModelState) {
	const minFraction = 1e-63
	selectionRangeLimit := int(math.Ceil(float64(ke.modelArchive.Len()) * ke.returnToBaseIsolationFraction))

	//selectedModel := ke.modelArchive.SelectRandomIsolatedModel(selectionRangeLimit)
	selectedModel := ke.modelArchive.SelectRandomModel()

	// Is this actually a problem?  Suppapitnarm paper doesn't mention it.  Original CRP doesn't cater for it.
	// It just seems odd to trawl through the isolated entries only to return to the current if its isolated.
	if currentModelState.Encoding() == selectedModel.Encoding() {
		warningMessage := fmt.Sprintf("Randomly selected return-to-base isolated model is same as current [%s]",
			currentModelState.Encoding())
		ke.LogHandler().Warn(warningMessage)
	}

	ke.modelArchive.Decompress(selectedModel, ke.currentModel)

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
		WithAttribute("New Base Model Encoding", selectedModel.Encoding())

	ke.lastReturnedToBase = ke.currentIteration

	ke.NotifyObserversOfEvent(*event)
}

func (ke *Explorer) returnToBaseOutsideArchive(currentModelState *archive.CompressedModelState) {
	const minFraction = 1e-63
	selectionRangeLimit := int(math.Ceil(float64(ke.modelArchive.Len()) * ke.returnToBaseIsolationFraction))

	selectedModel := ke.currentModel.DeepClone()
	selectedModel.Initialise()
	selectedModelArchive := ke.modelArchive.Compress(selectedModel)

	for currentModelState.Encoding() == selectedModelArchive.Encoding() {
		selectedModel.Initialise()
		selectedModelArchive = ke.modelArchive.Compress(selectedModel)
	}

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
		WithAttribute("New Base Model Encoding", selectedModelArchive.Encoding())

	ke.lastReturnedToBase = ke.currentIteration

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
	clone.currentModel = ke.currentModel.DeepClone()
	clone.potentialModel = ke.currentModel.DeepClone()
	return &clone
}

func (ke *Explorer) Model() model.Model {
	return ke.currentModel
}

func (ke *Explorer) SetModel(model model.Model) {
	ke.currentModel = model
	ke.potentialModel = model.DeepClone()
}

func (ke *Explorer) TearDown() {
	ke.LogHandler().Debug(ke.scenarioId + ": Triggering tear-down of Solution Explorer")
	ke.currentModel.TearDown()
	ke.potentialModel.TearDown()
}

func (ke *Explorer) setAcceptanceProbability(probability float64) {
	ke.coolant.SetAcceptanceProbability(math.Min(explorer.Guaranteed, probability))
}

func (ke *Explorer) EventAttributes(eventType observer.EventType) attributes.Attributes {
	switch eventType {
	case observer.StartedAnnealing:
		return ke.baseAttributes.
			Replace(explorer.Temperature, ke.coolant.Temperature()).
			Replace(ArchiveSize, ke.modelArchive.Len()).
			Add(explorer.CoolingFactor, ke.coolant.CoolingFactor())
	case observer.StartedIteration:
		return ke.baseAttributes.
			Replace(explorer.Temperature, ke.coolant.Temperature()).
			Replace(ArchiveSize, ke.modelArchive.Len())
	case observer.FinishedAnnealing:
		return ke.baseAttributes.
			Replace(explorer.Temperature, ke.coolant.Temperature()).
			Replace(ArchiveSize, ke.modelArchive.Len()).
			Add(ModelArchive, ke.modelArchive)
	case observer.FinishedIteration:
		return ke.baseAttributes.
			Replace(explorer.Temperature, ke.coolant.Temperature()).
			Replace(ArchiveSize, ke.modelArchive.Len()).
			Add(LastReturnedToBase, ke.lastReturnedToBase)
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
