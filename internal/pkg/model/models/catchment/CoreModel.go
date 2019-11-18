// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	baseParameters "github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

const (
	PlanningUnitsTableName = "PlanningUnits"
	GulliesTableName       = "Gullies"
)

func NewCoreModel() *CoreModel {
	newModel := new(CoreModel)

	newModel.SetName("CatchmentModel")
	newModel.SetEventNotifier(loggers.DefaultTestingEventNotifier)

	newModel.parameters.Initialise()
	newModel.managementActions.Initialise()
	newModel.ContainedDecisionVariables.Initialise()

	return newModel
}

type CoreModel struct {
	name.NameContainer
	name.IdentifiableContainer
	observer.ContainedEventNotifier

	parameters parameters.Parameters

	managementActions action.ModelManagementActions

	planningUnitTable tables.CsvTable
	gulliesTable      tables.CsvTable

	variable.ContainedDecisionVariables

	oleFunctionWrapper threading.MainThreadFunctionWrapper
	inputDataSet       dataset.DataSet
}

func (m *CoreModel) WithName(name string) *CoreModel {
	m.SetName(name)
	return m
}

func (m *CoreModel) WithParameters(params baseParameters.Map) *CoreModel {
	m.SetParameters(params)
	return m
}

func (m *CoreModel) WithSourceDataSet(sourceDataSet dataset.DataSet) *CoreModel {
	m.inputDataSet = sourceDataSet
	return m
}

func (m *CoreModel) SetParameters(params baseParameters.Map) error {
	m.parameters.AssignAllUserValues(params)
	return m.parameters.ValidationErrors()
}

func (m *CoreModel) ParameterErrors() error {
	return m.parameters.ValidationErrors()
}

func (m *CoreModel) Initialise() {
	m.planningUnitTable = m.fetchPlanningUnitTable()
	m.gulliesTable = m.fetchGulliesTable()

	m.buildDecisionVariables()
	m.buildManagementActions()
}

func (m *CoreModel) RandomlyInitialiseActions() {
	for _, action := range m.managementActions.Actions() {
		m.managementActions.RandomlyInitialiseAction(action)
	}
}

func (m *CoreModel) fetchPlanningUnitTable() tables.CsvTable {
	planningUnitTable, tableError := m.inputDataSet.Table(PlanningUnitsTableName)
	if tableError != nil {
		panic(errors.New("Expected data set supplied to have a [" + PlanningUnitsTableName + "] table"))
	}

	csvPlanningUnitTable, tableIsCsvType := planningUnitTable.(tables.CsvTable)
	if !tableIsCsvType {
		panic(errors.New("Expected data set table [" + PlanningUnitsTableName + "] to be a CSV type"))
	}
	return csvPlanningUnitTable
}

func (m *CoreModel) fetchGulliesTable() tables.CsvTable {
	gulliesTable, tableError := m.inputDataSet.Table(GulliesTableName)
	if tableError != nil {
		panic(errors.New("Expected data set supplied to have a [" + GulliesTableName + "] table"))
	}

	csvGulliesTable, tableIsCsvType := gulliesTable.(tables.CsvTable)
	if !tableIsCsvType {
		panic(errors.New("Expected data set table [" + GulliesTableName + "] to be a CSV type"))
	}
	return csvGulliesTable
}

func (m *CoreModel) buildDecisionVariables() {
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

func (m *CoreModel) buildManagementActions() {
	sedimentLoad := m.ContainedDecisionVariables.Variable(variables.SedimentProductionVariableName)
	implementationCost := m.ContainedDecisionVariables.Variable(variables.ImplementationCostVariableName)

	riverBankRestorations := new(actions.RiverBankRestorationGroup).Initialise(m.planningUnitTable, m.parameters)
	for _, action := range riverBankRestorations.ManagementActions() {
		m.managementActions.Add(action)
		action.Subscribe(m, sedimentLoad, implementationCost)
	}

	gullyRestorations := new(actions.GullyRestorationGroup).Initialise(m.gulliesTable, m.parameters)
	for _, action := range gullyRestorations.ManagementActions() {
		m.managementActions.Add(action)
		action.Subscribe(m, sedimentLoad, implementationCost)
	}

	hillSlopeRestorations := new(actions.HillSlopeRestorationGroup).Initialise(m.planningUnitTable, m.parameters)
	for _, action := range hillSlopeRestorations.ManagementActions() {
		m.managementActions.Add(action)
		action.Subscribe(m, sedimentLoad, implementationCost)
	}
}

func (m *CoreModel) ManagementActions() []action.ManagementAction {
	return m.managementActions.Actions()
}

func (m *CoreModel) ActiveManagementActions() []action.ManagementAction {
	return m.managementActions.ActiveActions()
}

func (m *CoreModel) SetManagementAction(index int, value bool) {
	m.managementActions.SetActivation(index, value)
}

func (m *CoreModel) SetManagementActionUnobserved(index int, value bool) {
	m.managementActions.SetActivationUnobserved(index, value)
}

func (m *CoreModel) PlanningUnits() planningunit.Ids {
	_, rows := m.planningUnitTable.ColumnAndRowSize()
	planningUnits := make(planningunit.Ids, rows)

	for row := uint(0); row < rows; row++ {
		planningUnit := m.planningUnitTable.CellFloat64(0, row)
		planningUnitId := planningunit.Float64ToId(planningUnit)
		planningUnits[row] = planningUnitId
	}

	return planningUnits
}

func (m *CoreModel) AcceptChange() {
	m.note("Accepting Change")
	m.ContainedDecisionVariables.AcceptAll()
}

func (m *CoreModel) RevertChange() {
	m.note("Reverting Change")
	m.ContainedDecisionVariables.RejectAll()
	m.managementActions.ToggleLastActivationUnobserved()
}

func (m *CoreModel) TearDown() {
	// deliberately does nothing.
}

func (m *CoreModel) DoRandomChange() {
	m.TryRandomChange()
	m.AcceptChange()
}

func (m *CoreModel) TryRandomChange() {
	m.note("Trying Random Change")
	m.managementActions.RandomlyToggleOneActivation()
}

func (m *CoreModel) UndoChange() {
	m.note("Undoing Change")
	m.managementActions.ToggleLastActivation()
}

func (m *CoreModel) ObserveAction(action action.ManagementAction) {
	m.noteAppliedManagementAction(action)
}

func (m *CoreModel) ObserveActionInitialising(action action.ManagementAction) {
	m.noteAppliedManagementAction(action)
}

func (m *CoreModel) noteAppliedManagementAction(action action.ManagementAction) {
	event := observer.NewEvent(observer.ManagementAction).
		WithId(m.Id()).
		WithAttribute("Type", action.Type()).
		WithAttribute("PlanningUnit", action.PlanningUnit())
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *CoreModel) note(text string) {
	event := observer.NewEvent(observer.Note).WithId(m.Id()).WithNote(text)
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *CoreModel) ObserveDecisionVariable(variable variable.DecisionVariable) {
	event := observer.NewEvent(observer.DecisionVariable).
		WithId(m.Id()).
		WithAttribute("Name", variable.Name()).
		WithAttribute("Value", variable.Value())
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *CoreModel) ObserveDecisionVariableWithNote(variable variable.DecisionVariable, note string) {
	event := observer.NewEvent(observer.DecisionVariable).
		WithId(m.Id()).
		WithAttribute("Name", variable.Name()).
		WithAttribute("Value", variable.Value()).
		WithNote(note)
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *CoreModel) capChangeOverRange(value float64) float64 {
	return math.Max(0, value)
}

func (m *CoreModel) DeepClone() model.Model {
	clone := *m
	clone.managementActions.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
