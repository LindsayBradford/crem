// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"math"
	"os"
	"path/filepath"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

const (
	downward = -1
	upward   = 1
)

const (
	PlanningUnitsTableName = "PlanningUnits"
)

func NewModel() *Model {
	newModel := new(Model)
	newModel.SetName("CatchmentModel")

	newModel.parameters.Initialise()

	return newModel
}

type Model struct {
	name.ContainedName
	rand.ContainedRand

	parameters Parameters

	ContainedDecisionVariables

	sedimentLoad *SedimentLoad

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

	planningUnitTable, tableError := m.dataSet.Table(PlanningUnitsTableName)
	if tableError != nil {
		panic(errors.New("Expected data set supplied to have a [" + PlanningUnitsTableName + "] table"))
	}

	csvPlanningUnitTable, ok := planningUnitTable.(*tables.CsvTable)
	if !ok {
		panic(errors.New("Expected data set table [" + PlanningUnitsTableName + "] to be a CSV type"))
	}

	sedimentLoad := new(SedimentLoad).Initialise(csvPlanningUnitTable, m.parameters)

	// TODO: Create Management actions
	// TODO: Have sedimentLoad subscribe to it's managemnt actions
	// TODO: Randomise which management actions are active

	// TODO: Bit of a smell to this.  Why don't I like it?
	m.decisionVariables.Add(&sedimentLoad.VolatileDecisionVariable)
}

func (m *Model) AcceptChange() {
	m.sedimentLoad.Accept()
}

func (m *Model) RevertChange() {
	m.sedimentLoad.Revert()
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
	return m.sedimentLoad.Value()
}

func (m *Model) setObjectiveValue(value float64) {
	m.sedimentLoad.SetTemporaryValue(value)
}

func (m *Model) DeepClone() model.Model {
	clone := *m
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
