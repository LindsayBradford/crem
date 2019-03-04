// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"math"
	"os"
	"path/filepath"

	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/actions"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/parameters"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/variables"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
	baseParameters "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

const (
	PlanningUnitsTableName = "PlanningUnits"
)

func NewModel() *Model {
	newModel := new(Model)
	newModel.SetName("CatchmentModel")

	newModel.parameters.Initialise()
	newModel.managementActions.Initialise()
	newModel.DecisionVariables().Initialise()

	return newModel
}

type Model struct {
	name.ContainedName
	name.ContainedIdentifier
	observer.ContainedEventNotifier

	parameters parameters.Parameters

	managementActions action.ModelManagementActions
	variables.ContainedDecisionVariables

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

func (m *Model) WithParameters(params baseParameters.Map) *Model {
	m.SetParameters(params)
	return m
}

func (m *Model) SetParameters(params baseParameters.Map) error {
	m.parameters.Merge(params)

	return m.parameters.ValidationErrors()
}

func (m *Model) ParameterErrors() error {
	return m.parameters.ValidationErrors()
}

func (m *Model) Initialise() {
	m.note("Initialising")

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

	sedimentLoad := new(variables.SedimentLoad).Initialise(csvPlanningUnitTable, m.parameters)
	sedimentLoad.Subscribe(m)

	implementationCost := new(variables.ImplementationCost).Initialise(csvPlanningUnitTable, m.parameters)
	implementationCost.Subscribe(m)

	// TODO: Create other sediment management actions
	riverBankRestorations := new(actions.RiverBankRestorations).Initialise(csvPlanningUnitTable, m.parameters)
	for _, action := range riverBankRestorations.ManagementActions() {
		m.managementActions.Add(action)
		action.Subscribe(m, sedimentLoad, implementationCost)
	}

	m.DecisionVariables().Add(
		&sedimentLoad.InductiveDecisionVariable,
		&implementationCost.InductiveDecisionVariable,
	)

	m.managementActions.RandomlyInitialise()
}

func (m *Model) AcceptChange() {
	m.note("Accepting Change")
	m.DecisionVariables().Accept()
}

func (m *Model) RevertChange() {
	m.note("Reverting Change")
	m.DecisionVariables().Revert()
	m.managementActions.UndoLastActivationToggleUnobserved()
}

func (m *Model) deriveDataSourcePath() string {
	relativeFilePath := m.parameters.GetString(parameters.DataSourcePath)
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, relativeFilePath)
}

func (m *Model) TearDown() {
	//  TODO: Do I need to do any special shutdown behaviour?
}

func (m *Model) TryRandomChange() {
	m.note("Trying Random Change")
	m.managementActions.RandomlyToggleOneActivation()
}

func (m *Model) ObserveAction(action action.ManagementAction) {
	m.noteAppliedManagementAction(action)
}

func (m *Model) ObserveActionInitialising(action action.ManagementAction) {
	m.noteAppliedManagementAction(action)
}

func (m *Model) noteAppliedManagementAction(action action.ManagementAction) {
	event := observer.NewEvent(observer.ManagementAction).
		WithId(m.Id()).
		WithSource(action)
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *Model) note(text string) {
	event := observer.NewEvent(observer.Note).
		WithId(m.Id()).
		WithSource(m).
		WithNote(text)
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *Model) ObserveDecisionVariable(variable variable.DecisionVariable) {
	event := observer.NewEvent(observer.DecisionVariable).
		WithId(m.Id()).
		WithSource(variable)
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *Model) capChangeOverRange(value float64) float64 {
	return math.Max(0, value)
}

func (m *Model) DeepClone() model.Model {
	clone := *m
	clone.managementActions.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
