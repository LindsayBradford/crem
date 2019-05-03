// Copyright (c) 2018 Australian Rivers Institute.

package kirkpatrick

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/pkg/errors"
)

type Explorer struct {
	name.NameContainer
	name.IdentifiableContainer

	model.ContainedModel
	rand.RandContainer
	loggers.ContainedLogger

	scenarioId string

	parameters            Parameters
	optimisationDirection optimisationDirection

	acceptanceProbability float64
	changeIsDesirable     bool
	changeAccepted        bool
	objectiveValueChange  float64

	temperature   float64
	coolingFactor float64
}

func New() *Explorer {
	newExplorer := new(Explorer)
	newExplorer.parameters.Initialise()
	newExplorer.SetModel(model.NewNullModel())
	return newExplorer
}

func (ke *Explorer) Initialise() {
	ke.LogHandler().Debug(ke.scenarioId + ": Initialising Solution Explorer")
	ke.SetRandomNumberGenerator(rand.NewTimeSeeded())
	ke.Model().Initialise()
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
	ke.parameters.Merge(params)

	ke.setOptimisationDirectionFromParams()

	ke.SetTemperature(ke.parameters.GetFloat64(StartingTemperature))
	ke.coolingFactor = ke.parameters.GetFloat64(CoolingFactor)

	ke.checkDecisionVariableFromParams()

	return ke
}

func (ke *Explorer) SetTemperature(temperature float64) error {
	if temperature <= 0 {
		return errors.New("invalid attempt to set annealer temperature to value <= 0")
	}
	ke.temperature = temperature
	return nil
}

func (ke *Explorer) setOptimisationDirectionFromParams() {
	optimisationDirectionParam := ke.parameters.GetString(OptimisationDirection)
	ke.optimisationDirection, _ = parseOptimisationDirection(optimisationDirectionParam)
}

func (ke *Explorer) checkDecisionVariableFromParams() {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)

	defer func() {
		if r := recover(); r != nil {
			ke.parameters.AddValidationErrorMessage("decision variable [" + decisionVariableName + "] not recognised by model")
		}
	}()

	ke.Model().DecisionVariable(decisionVariableName)
}

func (ke *Explorer) ParameterErrors() error {
	return ke.parameters.ValidationErrors()
}

func (ke *Explorer) ObjectiveValue() float64 {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)
	variable := ke.Model().DecisionVariable(decisionVariableName)
	return variable.Value()
}

func (ke *Explorer) TryRandomChange() {
	ke.Model().TryRandomChange()
	ke.defaultAcceptOrRevertChange()
}

func (ke *Explorer) defaultAcceptOrRevertChange() {
	ke.AcceptOrRevertChange(ke.AcceptLastChange, ke.RevertLastChange)
}

func (ke *Explorer) AcceptOrRevertChange(acceptFunction func(), revertFunction func()) {
	if ke.ChangeTriedIsDesirable() {
		ke.setAcceptanceProbability(explorer.Guaranteed)
		acceptFunction()
	} else {
		absoluteChangeInObjectiveValue := math.Abs(ke.ChangeInObjectiveValue())
		probabilityToAcceptBadChange := math.Exp(-absoluteChangeInObjectiveValue / ke.temperature)
		ke.setAcceptanceProbability(probabilityToAcceptBadChange)

		randomValue := ke.RandomNumberGenerator().Float64Unitary()
		if probabilityToAcceptBadChange > randomValue {
			acceptFunction()
		} else {
			revertFunction()
		}
	}
}

func (ke *Explorer) ChangeTriedIsDesirable() bool {
	switch ke.optimisationDirection {
	case Minimising:
		ke.SetChangeIsDesirable(ke.changeInObjectiveValue() <= 0)
		return ke.ChangeIsDesirable()
	case Maximising:
		ke.SetChangeIsDesirable(ke.changeInObjectiveValue() >= 0)
		return ke.ChangeIsDesirable()
	}
	return false
}

func (ke *Explorer) changeInObjectiveValue() float64 {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)
	ke.objectiveValueChange = ke.Model().DecisionVariableChange(decisionVariableName)
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

func (ke *Explorer) Temperature() float64 {
	return ke.temperature
}

func (ke *Explorer) ChangeIsDesirable() bool {
	return ke.changeIsDesirable
}

func (ke *Explorer) SetChangeIsDesirable(changeIsDesirable bool) {
	ke.changeIsDesirable = changeIsDesirable
}

func (ke *Explorer) ChangeInObjectiveValue() float64 {
	return ke.objectiveValueChange
}

func (ke *Explorer) ChangeAccepted() bool {
	return ke.changeAccepted
}

func (ke *Explorer) AcceptanceProbability() float64 {
	return ke.acceptanceProbability
}

func (ke *Explorer) setAcceptanceProbability(probability float64) {
	ke.acceptanceProbability = math.Min(explorer.Guaranteed, probability)
}

func (ke *Explorer) AttributesForEventType(eventType observer.EventType) attributes.Attributes {
	baseAttributes := new(attributes.Attributes).
		Add("ObjectiveValue", ke.ObjectiveValue()).
		Add("Temperature", ke.temperature)

	switch eventType {
	case observer.StartedAnnealing:
		return baseAttributes.Add("CoolingFactor", ke.coolingFactor)
	case observer.StartedIteration, observer.FinishedAnnealing:
		return baseAttributes
	case observer.FinishedIteration:
		return baseAttributes.
			Add("ChangeInObjectiveValue", ke.objectiveValueChange).
			Add("ChangeIsDesirable", ke.changeIsDesirable).
			Add("AcceptanceProbability", ke.acceptanceProbability).
			Add("ChangeAccepted", ke.changeAccepted)
	}
	return nil
}

func (ke *Explorer) CoolDown() {
	ke.temperature *= ke.coolingFactor
}
