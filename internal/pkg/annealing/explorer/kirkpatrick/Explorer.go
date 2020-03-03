// Copyright (c) 2018 Australian Rivers Institute.

package kirkpatrick

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/cooling/coolants/kirkpatrick"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model"
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
	Solution               = "Solution"
	ObjectiveValue         = "ObjectiveValue"
	ChangeInObjectiveValue = "ChangeInObjectiveValue"
)

type Explorer struct {
	name.NameContainer
	name.IdentifiableContainer

	model.ContainedModel
	loggers.ContainedLogger

	kirkpatrick.Coolant

	scenarioId string

	parameters            Parameters
	optimisationDirection optimisationDirection
	objectiveVariableName string

	changeIsDesirable    bool
	changeAccepted       bool
	changeInvalid        bool
	reasonChangeInvalid  string
	objectiveValueChange float64

	observer.SynchronousAnnealingEventNotifier
}

func New() *Explorer {
	newExplorer := new(Explorer)
	newExplorer.parameters.Initialise()
	newExplorer.Coolant.Initialise()
	newExplorer.SetModel(model.NewNullModel())
	return newExplorer
}

func (ke *Explorer) Initialise() {
	ke.LogHandler().Debug(ke.scenarioId + ": Initialising Solution Explorer")

	ke.notifyInitialisation()

	ke.SetRandomNumberGenerator(rand.NewTimeSeeded())
	ke.Model().Initialise()
}

func (ke *Explorer) notifyInitialisation() {
	event := observer.NewEvent(observer.Explorer).
		WithNote("Initialising").
		WithAttribute("OptimisationDirection", ke.optimisationDirection).
		WithAttribute("ObjectiveVariable", ke.objectiveVariableName)

	ke.NotifyObserversOfEvent(*event)
}

func (ke *Explorer) WithName(name string) *Explorer {
	ke.SetName(name)
	return ke
}

func (ke *Explorer) WithModel(model model.Model) *Explorer {
	ke.SetModel(model)
	return ke
}

func (ke *Explorer) WithParameters(params parameters.Map) *Explorer {
	ke.SetParameters(params)
	return ke
}

func (ke *Explorer) SetParameters(params parameters.Map) error {
	ke.parameters.AssignOnlyEnforcedUserValues(params)
	ke.Coolant.WithParameters(params)

	ke.setOptimisationDirectionFromParams()
	ke.checkDecisionVariableFromParams()

	return ke.parameters.ValidationErrors()
}

func (ke *Explorer) SetTemperature(temperature float64) error {
	if temperature <= 0 {
		return errors.New("invalid attempt to set annealer temperature to value <= 0")
	}
	ke.Temperature = temperature
	return nil
}

func (ke *Explorer) setOptimisationDirectionFromParams() {
	optimisationDirectionParam := ke.parameters.GetString(OptimisationDirection)
	ke.optimisationDirection, _ = parseOptimisationDirection(optimisationDirectionParam)
}

func (ke *Explorer) checkDecisionVariableFromParams() {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)

	if !ke.Model().OffersDecisionVariable(decisionVariableName) {
		ke.parameters.AddValidationErrorMessage("decision variable [" + decisionVariableName + "] not recognised by model")
	}

	ke.objectiveVariableName = decisionVariableName
}

func (ke *Explorer) ParameterErrors() error {
	mergedErrors := errors2.New("Kirkpatrick Explorer Parameter Validation")

	mergedErrors.Add(ke.parameters.ValidationErrors())
	mergedErrors.Add(ke.Coolant.ParameterErrors())

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
	ke.Model().TryRandomChange()
	ke.defaultAcceptOrRevertChange()
}

func (ke *Explorer) defaultAcceptOrRevertChange() {
	ke.AcceptOrRevertChange(ke.AcceptLastChange, ke.RevertLastChange)
}

func (ke *Explorer) AcceptOrRevertChange(acceptChange func(), revertChange func()) {
	ke.changeInvalid = false
	if isValid, invalidationErrors := ke.Model().ChangeIsValid(); !isValid {
		ke.reportInvalidChange(invalidationErrors)
		revertChange()
		return
	}

	if ke.changeTriedIsDesirable() {
		ke.setAcceptanceProbability(explorer.Guaranteed)
		ke.notifyDesirableAcceptance()
		acceptChange()
	} else {
		if ke.DecideIfAcceptable(ke.objectiveValueChange) {
			ke.notifyUndesirableAcceptance()
			acceptChange()
		} else {
			ke.notifyUndesirableReversion()
			revertChange()
		}
	}
}

func (ke *Explorer) notifyDesirableAcceptance() {
	event := observer.NewEvent(observer.Explorer).
		WithNote("Accepting Desirable Change").
		WithAttribute(explorer.AcceptanceProbability, ke.AcceptanceProbability)

	ke.NotifyObserversOfEvent(*event)
}

func (ke *Explorer) notifyUndesirableAcceptance() {
	acceptEvent := observer.NewEvent(observer.Explorer).
		WithNote("Accepting Undesirable Change").
		WithAttribute(explorer.AcceptanceProbability, ke.AcceptanceProbability)

	ke.NotifyObserversOfEvent(*acceptEvent)
}

func (ke *Explorer) notifyUndesirableReversion() {
	revertEvent := observer.NewEvent(observer.Explorer).
		WithNote("Reverting Undesirable Change").
		WithAttribute(explorer.AcceptanceProbability, ke.AcceptanceProbability)

	ke.NotifyObserversOfEvent(*revertEvent)
}

func (ke *Explorer) reportInvalidChange(invalidationErrors *errors2.CompositeError) {
	ke.calculateChangeInObjectiveValue()
	ke.changeInvalid = true
	ke.reasonChangeInvalid = invalidationErrors.Error()

	ke.notifyInvalidity()
}

func (ke *Explorer) notifyInvalidity() {
	event := observer.NewEvent(observer.Explorer).
		WithNote("Invalid Change").
		WithAttribute(ChangeInObjectiveValue, ke.objectiveValueChange).
		WithAttribute(explorer.ReasonChangeInvalid, ke.reasonChangeInvalid)

	ke.NotifyObserversOfEvent(*event)
}

func (ke *Explorer) changeTriedIsDesirable() bool {
	changeIsDesirable := false
	switch ke.optimisationDirection {
	case Minimising:
		changeIsDesirable = ke.calculateChangeInObjectiveValue() < 0
	case Maximising:
		changeIsDesirable = ke.calculateChangeInObjectiveValue() > 0
	}
	ke.changeIsDesirable = changeIsDesirable

	ke.notifyDesirability()

	return ke.changeIsDesirable
}

func (ke *Explorer) notifyDesirability() {
	event := observer.NewEvent(observer.Explorer).
		WithAttribute(ChangeInObjectiveValue, ke.objectiveValueChange).
		WithAttribute(explorer.ChangeIsDesirable, ke.changeIsDesirable)

	ke.NotifyObserversOfEvent(*event)
}

func (ke *Explorer) calculateChangeInObjectiveValue() float64 {
	ke.objectiveValueChange = ke.Model().DecisionVariableChange(ke.objectiveVariableName)
	return ke.objectiveValueChange
}

func (ke *Explorer) AcceptLastChange() {
	ke.Model().AcceptChange()
	ke.changeAccepted = true
}

func (ke *Explorer) RevertLastChange() {
	ke.Model().RevertChange()
	ke.changeAccepted = false
}

func (ke *Explorer) DeepClone() explorer.Explorer {
	clone := *ke
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	modelClone := ke.Model().DeepClone()
	clone.SetModel(modelClone)
	return &clone
}

func (ke *Explorer) TearDown() {
	ke.LogHandler().Debug(ke.scenarioId + ": Triggering tear-down of Solution Explorer")
	ke.Model().TearDown()
}

func (ke *Explorer) setAcceptanceProbability(probability float64) {
	ke.AcceptanceProbability = math.Min(explorer.Guaranteed, probability)
}

func (ke *Explorer) EventAttributes(eventType observer.EventType) attributes.Attributes {
	baseAttributes := new(attributes.Attributes).
		Add(ObjectiveValue, ke.ObjectiveValue()).
		Add(explorer.Temperature, ke.Temperature)

	switch eventType {
	case observer.StartedAnnealing:
		return new(attributes.Attributes).
			Add(explorer.Temperature, ke.Temperature).
			Add(explorer.CoolingFactor, ke.CoolingFactor).
			Add(ObjectiveValue, ke.ObjectiveValue())
	case observer.StartedIteration:
		return new(attributes.Attributes).
			Add(explorer.Temperature, ke.Temperature).
			Add(ObjectiveValue, ke.ObjectiveValue())
	case observer.FinishedAnnealing:
		return new(attributes.Attributes).
			Add(explorer.Temperature, ke.Temperature).
			Add(ObjectiveValue, ke.ObjectiveValue()).
			Add(Solution, *ke.fetchFinalModelSolution())
	case observer.Explorer:
		return baseAttributes
	case observer.FinishedIteration:
		return new(attributes.Attributes).Add(ObjectiveValue, ke.ObjectiveValue())
	}
	return nil
}

func (ke *Explorer) fetchFinalModelSolution() *solution.Solution {
	return new(solution.SolutionBuilder).
		WithId(ke.Id()).
		ForModel(ke.Model()).
		Build()
}

func (ke *Explorer) newEvent(eventType observer.EventType) {
	event := observer.NewEvent(eventType).
		JoiningAttributes(ke.EventAttributes(eventType))
	ke.NotifyObserversOfEvent(*event)
}

func (ke *Explorer) note(note string) {
	noteEvent := observer.NewEvent(observer.Explorer).
		WithAttribute(observer.Note.String(), note)
	ke.NotifyObserversOfEvent(*noteEvent)
}

func (ke *Explorer) CoolDown() {
	ke.Coolant.CoolDown()
	ke.notifyCoolDown()
}

func (ke *Explorer) notifyCoolDown() {
	event := observer.NewEvent(observer.Explorer).
		WithNote("Cooling").
		WithAttribute(explorer.Temperature, ke.Temperature)

	ke.NotifyObserversOfEvent(*event)
}
