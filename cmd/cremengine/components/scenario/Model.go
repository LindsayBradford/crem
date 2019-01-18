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
)

const (
	downward = -1
	upward   = 1
)

func NewModel() *Model {
	newModel := new(Model)
	newModel.SetName("Model")

	newModel.decisionVariables.Initialise()
	newModel.parameters.Initialise()

	return newModel
}

type Model struct {
	name.ContainedName
	rand.ContainedRand

	parameters Parameters

	decisionVariables DecisionVariables

	dataSet *excel.DataSet
}

func (m *Model) WithName(name string) *Model {
	m.SetName(name)
	return m
}

func (m *Model) WithOleFunctionWrapper(wrapper threading.MainThreadFunctionWrapper) *Model {
	m.dataSet = excel.NewDataSet("CatchmentDataSet", wrapper)
	return m
}

func (m *Model) WithParameters(params parameters.Map) *Model {
	m.SetParameters(params)
	return m
}

func (m *Model) SetParameters(params parameters.Map) error {
	m.parameters.Merge(params)

	return m.parameters.ValidationErrors()
}

func (m *Model) ParameterErrors() error {
	return m.parameters.ValidationErrors()
}

func (m *Model) Initialise() {
	m.SetRandomNumberGenerator(rand.NewTimeSeeded())
	dataSourcePath := m.deriveDataSourcePath()

	m.dataSet.Load(dataSourcePath)

	m.decisionVariables.SetValue(SedimentLoad, m.deriveInitialSedimentLoad())
}

func (m *Model) deriveInitialSedimentLoad() float64 {
	return 0
}

func (m *Model) DecisionVariable(name string) model.DecisionVariable {
	return m.decisionVariables.Variable(name)
}

func (m *Model) DecisionVariableChange(variableName string) float64 {
	decisionVariable := m.decisionVariables.Variable(variableName)
	difference := decisionVariable.TemporaryValue() - decisionVariable.Value()
	return difference
}

func (m *Model) AcceptChange() {
	m.decisionVariables.Variable(SedimentLoad).Accept()
}

func (m *Model) RevertChange() {
	m.decisionVariables.Variable(SedimentLoad).Revert()
}

func (m *Model) deriveDataSourcePath() string {
	relativeFilePath := m.parameters.GetString(DataSourcePath)
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, relativeFilePath)
}

func (m *Model) TearDown() {
	//  TODO: Do I need to do any special shutdown behaviour?
}

func (m *Model) TryRandomChange() {
	// TODO: randomly choose a management action to toggle.
	originalValue := m.objectiveValue()
	change := m.generateRandomChange()
	newValue := m.capChangeOverRange(originalValue + change)
	m.setObjectiveValue(newValue)
}

func (m *Model) generateRandomChange() float64 {
	randomValue := m.RandomNumberGenerator().Intn(2)

	var changeInObjectiveValue float64
	switch randomValue {
	case 0:
		changeInObjectiveValue = downward
	case 1:
		changeInObjectiveValue = upward
	}

	return changeInObjectiveValue
}

func (m *Model) capChangeOverRange(value float64) float64 {
	return math.Max(0, value)
}

func (m *Model) objectiveValue() float64 {
	return m.decisionVariables.Value(SedimentLoad)
}

func (m *Model) setObjectiveValue(value float64) {
	m.decisionVariables.Variable(SedimentLoad).SetTemporaryValue(value)
}

func (m *Model) DeepClone() model.Model {
	clone := *m
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
