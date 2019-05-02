// Copyright (c) 2018 Australian Rivers Institute.

package kirkpatrick

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/name"
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
	temperature           float64
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

func (ke *Explorer) WithId(id string) *Explorer {
	ke.SetId(id)
	return ke
}

func (ke *Explorer) WithParameters(params parameters.Map) *Explorer {
	ke.parameters.Merge(params)

	ke.setOptimisationDirectionFromParams()
	ke.checkDecisionVariableFromParams()

	return ke
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

func (ke *Explorer) TryRandomChange(temperature float64) {
	ke.Model().TryRandomChange()
	ke.defaultAcceptOrRevertChange(temperature)
}

func (ke *Explorer) defaultAcceptOrRevertChange(annealingTemperature float64) {
	ke.AcceptOrRevertChange(annealingTemperature, ke.AcceptLastChange, ke.RevertLastChange)
}

func (ke *Explorer) AcceptOrRevertChange(annealingTemperature float64, acceptFunction func(), revertFunction func()) {
	ke.SetTemperature(annealingTemperature)
	if ke.ChangeTriedIsDesirable() {
		ke.SetAcceptanceProbability(explorer.Guaranteed)
		acceptFunction()
	} else {
		absoluteChangeInObjectiveValue := math.Abs(ke.ChangeInObjectiveValue())
		probabilityToAcceptBadChange := math.Exp(-absoluteChangeInObjectiveValue / annealingTemperature)
		ke.SetAcceptanceProbability(probabilityToAcceptBadChange)

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
	change := ke.Model().DecisionVariableChange(decisionVariableName)
	ke.SetChangeInObjectiveValue(change)
	return change
}

func (ke *Explorer) AcceptLastChange() {
	ke.Model().AcceptChange()
	ke.SetChangeAccepted(true)
}

func (ke *Explorer) RevertLastChange() {
	ke.Model().RevertChange()
	ke.SetChangeAccepted(false)
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

func (ke *Explorer) SetTemperature(temperature float64) {
	ke.temperature = temperature
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

func (ke *Explorer) SetChangeInObjectiveValue(change float64) {
	ke.objectiveValueChange = change
}

func (ke *Explorer) ChangeAccepted() bool {
	return ke.changeAccepted
}

func (ke *Explorer) SetChangeAccepted(changeAccepted bool) {
	ke.changeAccepted = changeAccepted
}

func (ke *Explorer) AcceptanceProbability() float64 {
	return ke.acceptanceProbability
}

func (ke *Explorer) SetAcceptanceProbability(probability float64) {
	ke.acceptanceProbability = math.Min(explorer.Guaranteed, probability)
}
