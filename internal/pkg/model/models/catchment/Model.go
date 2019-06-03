// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	"math"
	"os"
	"path/filepath"
	"strconv"

	baseParameters "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/LindsayBradford/crem/pkg/strings"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

const (
	PlanningUnitsTableName = "PlanningUnits"
	GulliesTableName       = "Gullies"
)

func NewModel() *Model {
	newModel := new(Model)
	newModel.SetName("CatchmentModel")

	newModel.parameters.Initialise()
	newModel.managementActions.Initialise()
	newModel.ContainedDecisionVariables.Initialise()

	return newModel
}

type Model struct {
	name.NameContainer
	name.IdentifiableContainer
	observer.ContainedEventNotifier

	parameters parameters.Parameters

	managementActions action.ModelManagementActions

	planningUnitTable tables.CsvTable
	gulliesTable      tables.CsvTable

	variable.ContainedDecisionVariables

	oleFunctionWrapper threading.MainThreadFunctionWrapper
	inputDataSet       excel.DataSet
}

func (m *Model) WithName(name string) *Model {
	m.SetName(name)
	return m
}

func (m *Model) WithOleFunctionWrapper(wrapper threading.MainThreadFunctionWrapper) *Model {
	m.oleFunctionWrapper = wrapper
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

	m.inputDataSet = *excel.NewDataSet("CatchmentDataSet", m.oleFunctionWrapper)

	dataSourcePath := m.deriveDataSourcePath()

	m.inputDataSet.Load(dataSourcePath)

	m.planningUnitTable = m.fetchPlanningUnitTable()
	m.gulliesTable = m.fetchGulliesTable()

	m.buildCoreDecisionVariables()
	m.buildManagementActions()

	m.buildSedimentVsCostDecisionVariable()

	m.managementActions.RandomlyInitialise()
}

func (m *Model) fetchPlanningUnitTable() tables.CsvTable {

	planningUnitTable, tableError := m.inputDataSet.Table(PlanningUnitsTableName)
	if tableError != nil {
		panic(errors.New("Expected data set supplied to have a [" + PlanningUnitsTableName + "] table"))
	}

	csvPlanningUnitTable, ok := planningUnitTable.(tables.CsvTable)
	if !ok {
		panic(errors.New("Expected data set table [" + PlanningUnitsTableName + "] to be a CSV type"))
	}
	return csvPlanningUnitTable
}

func (m *Model) fetchGulliesTable() tables.CsvTable {

	gulliesTable, tableError := m.inputDataSet.Table(GulliesTableName)
	if tableError != nil {
		panic(errors.New("Expected data set supplied to have a [" + GulliesTableName + "] table"))
	}

	csvGulliesTable, ok := gulliesTable.(tables.CsvTable)
	if !ok {
		panic(errors.New("Expected data set table [" + GulliesTableName + "] to be a CSV type"))
	}
	return csvGulliesTable
}

func (m *Model) buildCoreDecisionVariables() {
	sedimentLoad := new(variables.SedimentProduction).
		Initialise(m.planningUnitTable, m.gulliesTable, m.parameters).
		WithObservers(m)

	implementationCost := new(variables.ImplementationCost).
		Initialise(m.planningUnitTable, m.parameters).
		WithObservers(m)

	m.ContainedDecisionVariables.Add(
		sedimentLoad,
		implementationCost,
	)
}

func (m *Model) buildManagementActions() {
	sedimentLoad := m.ContainedDecisionVariables.Variable(variables.SedimentProductionVariableName)
	implementationCost := m.ContainedDecisionVariables.Variable(variables.ImplementationCostVariableName)

	// TODO: Create other sediment management actions
	riverBankRestorations := new(actions.RiverBankRestorations).Initialise(m.planningUnitTable, m.parameters)
	for _, action := range riverBankRestorations.ManagementActions() {
		m.managementActions.Add(action)
		action.Subscribe(m, sedimentLoad, implementationCost)
	}

	gullyRestorations := new(actions.GullyRestorations).Initialise(m.gulliesTable, m.parameters)
	for _, action := range gullyRestorations.ManagementActions() {
		m.managementActions.Add(action)
		action.Subscribe(m, sedimentLoad, implementationCost)
	}
}

func (m *Model) ManagementActions() []action.ManagementAction {
	return m.managementActions.Actions()
}

func (m *Model) ActiveManagementActions() []action.ManagementAction {
	return m.managementActions.ActiveActions()
}

func (m *Model) SetManagementAction(index int, value bool) {
	m.managementActions.SetActivation(index, value)
}

func (m *Model) PlanningUnits() solution.PlanningUnitIds {
	_, rows := m.planningUnitTable.ColumnAndRowSize()
	planningUnits := make(solution.PlanningUnitIds, rows)

	for row := uint(0); row < rows; row++ {
		planningUnit := m.planningUnitTable.CellFloat64(0, row)
		planningUnitAsString := strconv.FormatFloat(planningUnit, 'g', -1, 64)
		planningUnitId := planningUnitAsString
		planningUnits[row] = planningUnitId
	}

	return planningUnits
}

func (m *Model) buildSedimentVsCostDecisionVariable() {
	sedimentLoad := m.ContainedDecisionVariables.Variable(variables.SedimentProductionVariableName)
	implementationCost := m.ContainedDecisionVariables.Variable(variables.ImplementationCostVariableName)

	const sedimentWeight = 0.5
	const implementationWeight = 0.5

	sedimentVsCost, buildError := new(variables.SedimentVsCost).
		Initialise().
		WithObservers(m).
		WithWeightedVariable(sedimentLoad, sedimentWeight).
		WithWeightedVariable(implementationCost, implementationWeight).
		Build()

	if buildError != nil {
		panic(buildError)
	}

	noteBuilder := new(strings.FluentBuilder).
		Add(sedimentLoad.Name(), " weight = ", strconv.FormatFloat(sedimentWeight, 'f', 3, 64), ", ").
		Add(implementationCost.Name(), " weight = ", strconv.FormatFloat(implementationWeight, 'f', 3, 64))

	m.ObserveDecisionVariableWithNote(sedimentVsCost, noteBuilder.String())

	m.ContainedDecisionVariables.Add(sedimentVsCost)
}

func (m *Model) AcceptChange() {
	m.note("Accepting Change")
	m.ContainedDecisionVariables.AcceptAll()
}

func (m *Model) RevertChange() {
	m.note("Reverting Change")
	m.ContainedDecisionVariables.RejectAll()
	m.managementActions.UndoLastActivationToggleUnobserved()
}

func (m *Model) deriveDataSourcePath() string {
	relativeFilePath := m.parameters.GetString(parameters.DataSourcePath)
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, relativeFilePath)
}

func (m *Model) TearDown() {
	m.inputDataSet.Teardown()
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
		WithAttribute("Type", action.Type()).
		WithAttribute("PlanningUnit", action.PlanningUnit()).
		WithAttribute("IsActive", action.IsActive())
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *Model) note(text string) {
	event := observer.NewEvent(observer.Note).WithId(m.Id()).WithNote(text)
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *Model) ObserveDecisionVariable(variable variable.DecisionVariable) {
	event := observer.NewEvent(observer.DecisionVariable).
		WithId(m.Id()).
		WithAttribute("Name", variable.Name()).
		WithAttribute("Value", variable.Value())
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *Model) ObserveDecisionVariableWithNote(variable variable.DecisionVariable, note string) {
	event := observer.NewEvent(observer.DecisionVariable).
		WithId(m.Id()).
		WithAttribute("Name", variable.Name()).
		WithAttribute("Value", variable.Value()).
		WithNote(note)
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
