// Copyright (c) 2018 Australian Rivers Institute.

package kirkpatrick

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

const guaranteed = 1

type Explorer struct {
	name.ContainedName
	model.ContainedModel
	rand.ContainedRand
	logging.ContainedLogger

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
	if modelWithScenarioId, hasScenarioId := ke.Model().(scenario.Identifiable); hasScenarioId {
		modelWithScenarioId.SetScenarioId(ke.ScenarioId())
	}
}

func (ke *Explorer) WithName(name string) *Explorer {
	ke.SetName(name)
	return ke
}

func (ke *Explorer) WithModel(model model.Model) *Explorer {
	ke.SetModel(model)
	return ke
}

func (ke *Explorer) ScenarioId() string {
	return ke.scenarioId
}

func (ke *Explorer) SetScenarioId(id string) {
	ke.scenarioId = id
}

func (ke *Explorer) WithScenarioId(id string) *Explorer {
	ke.scenarioId = id
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
	_, dvError := ke.Model().DecisionVariable(decisionVariableName)
	if dvError != nil {
		ke.parameters.AddValidationErrorMessage("decision variable [" + decisionVariableName + "] not recognised by model")
	}
}

func (ke *Explorer) ParameterErrors() error {
	return ke.parameters.ValidationErrors()
}

func (ke *Explorer) ObjectiveValue() float64 {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)
	if dv, dvError := ke.Model().DecisionVariable(decisionVariableName); dvError == nil {
		return dv.Value()
	}
	return 0
}

func (ke *Explorer) TryRandomChange(temperature float64) {
	ke.Model().TryRandomChange()
	ke.defaultAcceptOrRevertChange(temperature)
}

func (ke *Explorer) defaultAcceptOrRevertChange(annealingTemperature float64) {
	ke.AcceptOrRevertChange(annealingTemperature, ke.AcceptLastChange, ke.RevertLastChange)
}

func (ke *Explorer) AcceptOrRevertChange(annealingTemperature float64, acceptFunction func(), revertFunction func()) {
	ke.temperature = annealingTemperature
	if ke.ChangeTriedIsDesirable() {
		ke.SetAcceptanceProbability(guaranteed)
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
		ke.changeIsDesirable = ke.changeInObjectiveValue() < 0
		return ke.changeIsDesirable
	case Maximising:
		ke.changeIsDesirable = ke.changeInObjectiveValue() > 0
		return ke.changeIsDesirable
	}
	return false
}

func (ke *Explorer) Temperature() float64 {
	return ke.temperature
}

func (ke *Explorer) ChangeIsDesirable() bool {
	return ke.changeIsDesirable
}

func (ke *Explorer) changeInObjectiveValue() float64 {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)
	if change, changeError := ke.Model().DecisionVariableChange(decisionVariableName); changeError == nil {
		ke.SetChangeInObjectiveValue(change)
		return change
	}
	return 0
}

func (ke *Explorer) ChangeInObjectiveValue() float64 {
	return ke.objectiveValueChange
}

func (ke *Explorer) SetChangeInObjectiveValue(change float64) {
	ke.objectiveValueChange = change
}

func (ke *Explorer) AcceptLastChange() {
	ke.Model().AcceptChange()
	ke.changeAccepted = true
}

func (ke *Explorer) RevertLastChange() {
	ke.Model().RevertChange()
	ke.changeAccepted = false
}

func (ke *Explorer) ChangeAccepted() bool {
	return ke.changeAccepted
}

func (ke *Explorer) AcceptanceProbability() float64 {
	return ke.acceptanceProbability
}

func (ke *Explorer) SetAcceptanceProbability(probability float64) {
	ke.acceptanceProbability = math.Min(guaranteed, probability)
}

func (ke *Explorer) DeepClone() explorer.Explorer {
	clone := *ke
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	modelClone := ke.Model().DeepClone()
	clone.SetModel(modelClone)
	return &clone
}

func (ke *Explorer) CloneObservable() explorer.Explorer {
	clone := *ke
	clone.SetModel(model.NullModel)
	return &clone
}

func (ke *Explorer) TearDown() {
	ke.LogHandler().Debug(ke.scenarioId + ": Triggering tear-down of Solution Explorer")
	ke.Model().TearDown()
}
