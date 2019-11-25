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
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	baseParameters "github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/name"
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

	inputDataSet dataset.DataSet
	initialising bool
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
	m.planningUnitTable = m.fetchCsvTable(PlanningUnitsTableName)
	m.gulliesTable = m.fetchCsvTable(GulliesTableName)

	m.buildDecisionVariables()
	m.buildAndObserveManagementActions()
}

func (m *CoreModel) fetchCsvTable(tableName string) tables.CsvTable {
	namedTable, namedTableError := m.inputDataSet.Table(tableName)
	if namedTableError != nil {
		panic(errors.New("Expected data set supplied to have a [" + tableName + "] table"))
	}

	namedCsvTable, isCsvType := namedTable.(tables.CsvTable)
	if !isCsvType {
		panic(errors.New("Expected data set table [" + tableName + "] to be a CSV type"))
	}
	return namedCsvTable
}

func (m *CoreModel) buildDecisionVariables() {
	sedimentProduction := new(variables.SedimentProduction). // TODO: retire this when sedimentProduction2 finalised.
									Initialise(m.planningUnitTable, m.gulliesTable, m.parameters).
									WithObservers(m)

	sedimentProduction2 := new(variables.SedimentProduction2).
		Initialise(m.planningUnitTable, m.gulliesTable, m.parameters).
		WithObservers(m)

	implementationCost := new(variables.ImplementationCost). // TODO: retire this when implementationCost2 finalised.
									Initialise(m.planningUnitTable, m.parameters).
									WithObservers(m)

	implementationCost2 := new(variables.ImplementationCost2).
		Initialise(m.planningUnitTable, m.parameters).
		WithObservers(m)

	m.ContainedDecisionVariables.Add(
		sedimentProduction, sedimentProduction2,
		implementationCost, implementationCost2,
	)
}

func (m *CoreModel) buildAndObserveManagementActions() {
	actions := m.buildModelActions()
	observers := m.buildActionObservers()
	m.observeActions(observers, actions)
}

func (m *CoreModel) buildModelActions() []action.ManagementAction {
	modelActions := make([]action.ManagementAction, 0)

	modelActions = append(modelActions, m.buildRiverBankRestorations()...)
	modelActions = append(modelActions, m.buildGullyRestorations()...)
	modelActions = append(modelActions, m.buildHillSlopeRestorations()...)

	return modelActions
}

func (m *CoreModel) buildRiverBankRestorations() []action.ManagementAction {
	riverBankRestorations := new(actions.RiverBankRestorationGroup).
		WithPlanningUnitTable(m.planningUnitTable).
		WithParameters(m.parameters).
		ManagementActions()
	return riverBankRestorations
}

func (m *CoreModel) buildGullyRestorations() []action.ManagementAction {
	gullyRestorations := new(actions.GullyRestorationGroup).
		WithParameters(m.parameters).
		WithGullyTable(m.gulliesTable).
		ManagementActions()
	return gullyRestorations
}

func (m *CoreModel) buildHillSlopeRestorations() []action.ManagementAction {
	hillSlopeRestorations := new(actions.HillSlopeRestorationGroup).
		WithPlanningUnitTable(m.planningUnitTable).
		WithParameters(m.parameters).
		ManagementActions()
	return hillSlopeRestorations
}

func (m *CoreModel) buildActionObservers() []action.Observer {
	sedimentProduction := m.ContainedDecisionVariables.Variable(variables.SedimentProductionVariableName)
	sedimentProduction2 := m.ContainedDecisionVariables.Variable(variables.SedimentProduction2VariableName)

	implementationCost := m.ContainedDecisionVariables.Variable(variables.ImplementationCostVariableName)
	implementationCost2 := m.ContainedDecisionVariables.Variable(variables.ImplementationCost2VariableName)

	return []action.Observer{
		m, sedimentProduction, sedimentProduction2, implementationCost, implementationCost2,
	}
}

func (m *CoreModel) observeActions(actionObservers []action.Observer, actions []action.ManagementAction) {
	for _, action := range actions {
		m.managementActions.Add(action)
		action.Subscribe(actionObservers...)
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
	if m.initialising {
		return
	}
	m.note("Accepting Change")
	m.ContainedDecisionVariables.AcceptAll()
}

func (m *CoreModel) RevertChange() {
	m.note("Reverting Change")
	m.ContainedDecisionVariables.RejectAll()
	m.managementActions.ToggleLastActivationUnobserved()
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
	if m.initialising {
		return
	}
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

func (m *CoreModel) ObserveDecisionVariable(variable variableNew.DecisionVariable) {
	if m.initialising {
		return
	}
	event := observer.NewEvent(observer.DecisionVariable).
		WithId(m.Id()).
		WithAttribute("Name", variable.Name()).
		WithAttribute("Value", variable.Value())
	m.EventNotifier().NotifyObserversOfEvent(*event)
}

func (m *CoreModel) ObserveDecisionVariableWithNote(variable variableNew.DecisionVariable, note string) {
	if m.initialising {
		return
	}
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

func (m *CoreModel) TearDown() {
	// deliberately does nothing.
}
