// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"errors"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/dumb"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
)

const (
	_                                  = iota
	BankErosionFudgeFactor      string = "BankErosionFudgeFactor"
	WaterDensity                string = "WaterDensity"
	LocalAcceleration           string = "LocalAcceleration"
	GullyCompensationFactor     string = "GullyCompensationFactor"
	SedimentDensity             string = "SedimentDensity"
	SuspendedSedimentProportion string = "SuspendedSedimentProportion"
)

type CatchmentModel struct {
	name.ContainedName
	rand.ContainedRand

	decisionVariables     map[string]model.DecisionVariable
	tempDecisionVariables map[string]model.DecisionVariable

	parameters CatchmentParameters
}

func NewCatchmentModel() *CatchmentModel {
	newModel := new(CatchmentModel)
	newModel.SetName("CatchmentModel")

	newModel.SetRandomNumberGenerator(rand.NewTimeSeeded())
	newModel.buildDecisionVariables()
	newModel.parameters.Initialise()

	return newModel
}

func (cm *CatchmentModel) buildDecisionVariables() {
	cm.decisionVariables = make(map[string]model.DecisionVariable, 1)
	cm.tempDecisionVariables = make(map[string]model.DecisionVariable, 1)

	objectiveValueVar := new(model.DecisionVariableImpl)
	objectiveValueVar.SetName(model.ObjectiveValue)
	cm.decisionVariables[model.ObjectiveValue] = objectiveValueVar

	tempObjectiveValueVar := new(model.DecisionVariableImpl)
	tempObjectiveValueVar.SetName(model.ObjectiveValue)
	cm.tempDecisionVariables[model.ObjectiveValue] = tempObjectiveValueVar
}

func (cm *CatchmentModel) WithName(name string) *CatchmentModel {
	cm.SetName(name)
	return cm
}

func (cm *CatchmentModel) WithParameters(params parameters.Map) *CatchmentModel {
	cm.SetParameters(params)
	return cm
}

func (cm *CatchmentModel) SetParameters(params parameters.Map) error {
	cm.parameters.Merge(params)

	initialValue := cm.parameters.GetFloat64(dumb.InitialObjectiveValue)
	cm.decisionVariables[model.ObjectiveValue].SetValue(initialValue)
	cm.tempDecisionVariables[model.ObjectiveValue].SetValue(initialValue)

	return cm.parameters.ValidationErrors()
}

func (cm *CatchmentModel) ParameterErrors() error {
	return cm.parameters.ValidationErrors()
}

func (cm *CatchmentModel) Initialise() {
	// This model doesn't need any special initialising.
}

func (cm *CatchmentModel) TearDown() {
	// This model doesn't need any special tearDown.
}

func (cm *CatchmentModel) TryRandomChange() {
	// TODO: randomly choose a management action to toggle.
}

func (cm *CatchmentModel) copyDecisionVarValueToTemp(varName string) {
	cm.tempDecisionVariables[varName].SetValue(
		cm.decisionVariables[varName].Value(),
	)
}

func (cm *CatchmentModel) copyTempDecisionVarValueToActual(varName string) {
	cm.decisionVariables[varName].SetValue(
		cm.tempDecisionVariables[varName].Value(),
	)
}

func (cm *CatchmentModel) AcceptChange() {
	cm.copyTempDecisionVarValueToActual(model.ObjectiveValue)
}

func (cm *CatchmentModel) RevertChange() {
	cm.copyDecisionVarValueToTemp(model.ObjectiveValue)
}

func (cm *CatchmentModel) DecisionVariable(name string) (model.DecisionVariable, error) {
	if variable, found := cm.decisionVariables[name]; found == true {
		return variable, nil
	}
	return model.NullDecisionVariable, errors.New("decision variable [" + name + "] not defined for model [" + cm.Name() + " ].")
}

func (cm *CatchmentModel) DecisionVariableChange(variableName string) (float64, error) {
	decisionVariable, foundActual := cm.decisionVariables[variableName]
	tmpDecisionVar, foundTemp := cm.tempDecisionVariables[decisionVariable.Name()]
	if !foundActual || !foundTemp {
		return 0, errors.New("no temporary decision variable of name [" + decisionVariable.Name() + "] in model [" + cm.Name() + "].")
	}

	difference := tmpDecisionVar.Value() - decisionVariable.Value()
	return difference, nil
}

func (cm *CatchmentModel) DeepClone() model.Model {
	clone := *cm
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
