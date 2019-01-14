// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"math"
	"os"
	"path/filepath"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

const (
	downward = -1
	upward   = 1
)

type CatchmentModel struct {
	name.ContainedName
	rand.ContainedRand

	parameters CatchmentParameters

	decisionVariables model.VolatileDecisionVariables

	dataSet *excel.DataSet
}

func NewCatchmentModel() *CatchmentModel {
	newModel := new(CatchmentModel)
	newModel.SetName("CatchmentModel")

	newModel.BuildDecisionVariables()
	newModel.parameters.Initialise()

	return newModel
}

func (cm *CatchmentModel) BuildDecisionVariables() {
	cm.createDecisionVariables()
}

func (cm *CatchmentModel) createDecisionVariables() {
	cm.decisionVariables = model.NewVolatileDecisionVariables()
	cm.createDecisionVariableFor(SedimentLoad)
}

func (cm *CatchmentModel) createDecisionVariableFor(name string) {
	newVariable := new(model.VolatileDecisionVariable)
	newVariable.SetName(name)
	cm.decisionVariables[name] = newVariable
}

func (cm *CatchmentModel) WithName(name string) *CatchmentModel {
	cm.SetName(name)
	return cm
}

func (cm *CatchmentModel) WithOleFunctionWrapper(wrapper threading.MainThreadFunctionWrapper) *CatchmentModel {
	cm.dataSet = excel.NewDataSet("CatchmentDataSet", wrapper)
	return cm
}

func (cm *CatchmentModel) WithParameters(params parameters.Map) *CatchmentModel {
	cm.SetParameters(params)
	return cm
}

func (cm *CatchmentModel) SetParameters(params parameters.Map) error {
	cm.parameters.Merge(params)

	return cm.parameters.ValidationErrors()
}

func (cm *CatchmentModel) ParameterErrors() error {
	return cm.parameters.ValidationErrors()
}

func (cm *CatchmentModel) Initialise() {
	cm.SetRandomNumberGenerator(rand.NewTimeSeeded())
	dataSourcePath := cm.deriveDataSourcePath()

	cm.dataSet.Load(dataSourcePath)

	initialValue := cm.parameters.GetFloat64(SedimentLoad)
	cm.SetDecisionVariable(SedimentLoad, initialValue)
}

func (cm *CatchmentModel) SetDecisionVariable(name string, value float64) {
	cm.decisionVariables[name].SetValue(value)
}

func (cm *CatchmentModel) DecisionVariableChange(variableName string) (float64, error) {
	decisionVariable, foundActual := cm.decisionVariables[variableName]
	if !foundActual {
		return 0, errors.New("no temporary decision variable of name [" + decisionVariable.Name() + "] in model [" + cm.Name() + "].")
	}

	difference := decisionVariable.TemporaryValue() - decisionVariable.Value()

	return difference, nil
}

func (cm *CatchmentModel) AcceptChange() {
	cm.decisionVariables[SedimentLoad].Accept()
}

func (cm *CatchmentModel) RevertChange() {
	cm.decisionVariables[SedimentLoad].Revert()
}

func (cm *CatchmentModel) DecisionVariable(name string) (model.DecisionVariable, error) {
	if variable, found := cm.decisionVariables[name]; found == true {
		return variable, nil
	}
	return model.NullDecisionVariable, errors.New("decision variable [" + name + "] not defined for model [" + cm.Name() + " ].")
}

func (cm *CatchmentModel) deriveDataSourcePath() string {
	relativeFilePath := cm.parameters.GetString(DataSourcePath)
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, relativeFilePath)
}

func (cm *CatchmentModel) TearDown() {
	//  TODO: Do I need to do any special shutdown behaviour?
}

func (cm *CatchmentModel) TryRandomChange() {
	// TODO: randomly choose a management action to toggle.
	originalValue := cm.objectiveValue()
	change := cm.generateRandomChange()
	newValue := cm.capChangeOverRange(originalValue + change)
	cm.setObjectiveValue(newValue)
}

func (cm *CatchmentModel) generateRandomChange() float64 {
	randomValue := cm.RandomNumberGenerator().Intn(2)

	var changeInObjectiveValue float64
	switch randomValue {
	case 0:
		changeInObjectiveValue = downward
	case 1:
		changeInObjectiveValue = upward
	}

	return changeInObjectiveValue
}

func (cm *CatchmentModel) capChangeOverRange(value float64) float64 {
	return math.Max(0, value)
}

func (cm *CatchmentModel) objectiveValue() float64 {
	return cm.decisionVariables[SedimentLoad].Value()
}

func (cm *CatchmentModel) setObjectiveValue(value float64) {
	cm.decisionVariables[SedimentLoad].SetTemporaryValue(value)
}

func (cm *CatchmentModel) DeepClone() model.Model {
	clone := *cm
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
