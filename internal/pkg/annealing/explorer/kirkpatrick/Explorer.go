// Copyright (c) 2018 Australian Rivers Institute.

package kirkpatrick

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
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
	objectiveVariableName string

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
	ke.objectiveVariableName = decisionVariableName
}

func (ke *Explorer) ParameterErrors() error {
	return ke.parameters.ValidationErrors()
}

func (ke *Explorer) ObjectiveValue() float64 {
	variable := ke.Model().DecisionVariable(ke.objectiveVariableName)
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
	if ke.changeTriedIsDesirable() {
		ke.setAcceptanceProbability(explorer.Guaranteed)
		acceptFunction()
	} else {
		if ke.shouldAcceptBadChange() {
			acceptFunction()
		} else {
			revertFunction()
		}
	}
}

func (ke *Explorer) shouldAcceptBadChange() bool {
	probabilityToAcceptBadChange := ke.calculateProbabilityToAcceptBadChange()
	randomValue := ke.RandomNumberGenerator().Float64Unitary()
	return probabilityToAcceptBadChange > randomValue
}

func (ke *Explorer) calculateProbabilityToAcceptBadChange() float64 {
	absoluteChangeInObjectiveValue := math.Abs(ke.objectiveValueChange)
	probabilityToAcceptBadChange := math.Exp(-absoluteChangeInObjectiveValue / ke.temperature)

	ke.setAcceptanceProbability(probabilityToAcceptBadChange)

	return probabilityToAcceptBadChange
}

func (ke *Explorer) changeTriedIsDesirable() bool {
	switch ke.optimisationDirection {
	case Minimising:
		ke.changeIsDesirable = ke.calculateChangeInObjectiveValue() <= 0
	case Maximising:
		ke.changeIsDesirable = ke.calculateChangeInObjectiveValue() >= 0
	}
	return ke.changeIsDesirable
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
	ke.acceptanceProbability = math.Min(explorer.Guaranteed, probability)
}

func (ke *Explorer) EventAttributes(eventType observer.EventType) attributes.Attributes {
	baseAttributes := new(attributes.Attributes).
		Add(explorer.ObjectiveValue, ke.ObjectiveValue()).
		Add(explorer.Temperature, ke.temperature)

	switch eventType {
	case observer.StartedAnnealing:
		return baseAttributes.Add(explorer.CoolingFactor, ke.coolingFactor)
	case observer.StartedIteration, observer.FinishedAnnealing:
		return baseAttributes
	case observer.FinishedIteration:
		return baseAttributes.
			Add(explorer.ChangeInObjectiveValue, ke.objectiveValueChange).
			Add(explorer.ChangeIsDesirable, ke.changeIsDesirable).
			Add(explorer.AcceptanceProbability, ke.acceptanceProbability).
			Add(explorer.ChangeAccepted, ke.changeAccepted)
	}
	return nil
}

func (ke *Explorer) CoolDown() {
	ke.temperature *= ke.coolingFactor
}
