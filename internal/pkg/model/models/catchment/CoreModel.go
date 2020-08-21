// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/opportunitycost"
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/implementationcost"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/nitrogenproduction"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/sedimentproduction"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	baseParameters "github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/pkg/errors"
)

const (
	SubcatchmentsTableName = "Subcatchments"
	GulliesTableName       = "Gullies"
	actionsTableName       = "Actions"
)

func NewCoreModel() *CoreModel {
	newModel := new(CoreModel)

	newModel.SetName("CatchmentModel")
	newModel.SetId("TestCatchmentModel")

	newModel.parameters.Initialise()
	newModel.managementActions.Initialise()
	newModel.ContainedDecisionVariables.Initialise()

	return newModel
}

type CoreModel struct {
	name.NameContainer
	name.IdentifiableContainer
	parameters parameters.Parameters

	managementActions action.ModelManagementActions

	planningUnitTable tables.CsvTable
	gulliesTable      tables.CsvTable
	actionsTable      tables.CsvTable

	variable.ContainedDecisionVariables

	inputDataSet dataset.DataSet
	initialising bool

	observer.SynchronousAnnealingEventNotifier
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
	m.validateModelParameters()
	return m.parameters.ValidationErrors()
}

func (m *CoreModel) validateModelParameters() {
	boundVariableNumber := 0
	if m.parameters.HasEntry(parameters.MaximumImplementationCost) {
		boundVariableNumber++
	}
	if m.parameters.HasEntry(parameters.MaximumSedimentProduction) {
		boundVariableNumber++
	}
	if m.parameters.HasEntry(parameters.MaximumParticulateNitrogenProduction) {
		boundVariableNumber++
	}

	if boundVariableNumber > 1 {
		errorText := fmt.Sprintf("Only one of [%s], [%s] or [%s] allowed as variable limit.",
			parameters.MaximumImplementationCost,
			parameters.MaximumSedimentProduction,
			parameters.MaximumParticulateNitrogenProduction)

		m.parameters.AddValidationErrorMessage(errorText)
	}
}

func (m *CoreModel) ParameterErrors() error {
	return m.parameters.ValidationErrors()
}

func (m *CoreModel) Initialise() {
	m.planningUnitTable = m.fetchCsvTable(SubcatchmentsTableName)
	m.gulliesTable = m.fetchCsvTable(GulliesTableName)
	m.actionsTable = m.fetchCsvTable(actionsTableName)

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
	sedimentProduction := new(sedimentproduction.SedimentProduction).
		Initialise(m.planningUnitTable, m.gulliesTable, m.parameters).
		WithObservers(m)

	if m.parameters.HasEntry(parameters.MaximumSedimentProduction) {
		sedimentProduction.SetMaximum(m.parameters.GetFloat64(parameters.MaximumSedimentProduction))
	}

	nitrogenProduction := new(nitrogenproduction.ParticulateNitrogenProduction).
		Initialise(m.actionsTable, m.parameters).
		WithObservers(m)

	if m.parameters.HasEntry(parameters.MaximumParticulateNitrogenProduction) {
		nitrogenProduction.SetMaximum(m.parameters.GetFloat64(parameters.MaximumParticulateNitrogenProduction))
	}

	implementationCost := new(implementationcost.ImplementationCost).
		Initialise(m.planningUnitTable, m.parameters).
		WithObservers(m)

	if m.parameters.HasEntry(parameters.MaximumImplementationCost) {
		implementationCost.SetMaximum(m.parameters.GetFloat64(parameters.MaximumImplementationCost))
	}

	opportunityCost := new(opportunitycost.OpportunityCost).
		Initialise(m.planningUnitTable, m.parameters).
		WithObservers(m)

	if m.parameters.HasEntry(parameters.MaximumOpportunityCost) {
		opportunityCost.SetMaximum(m.parameters.GetFloat64(parameters.MaximumOpportunityCost))
	}

	m.ContainedDecisionVariables.Initialise()
	m.ContainedDecisionVariables.Add(
		sedimentProduction, nitrogenProduction, implementationCost, opportunityCost,
	)
}

func (m *CoreModel) buildAndObserveManagementActions() {
	actions := m.buildModelActions()
	observers := m.buildActionObservers()
	m.observeActions(observers, actions)
}

func (m *CoreModel) buildModelActions() []action.ManagementAction {
	modelActions := make([]action.ManagementAction, 0)

	modelActions = append(modelActions, m.buildGullyRestorations()...)
	modelActions = append(modelActions, m.buildRiverBankRestorations()...)
	modelActions = append(modelActions, m.buildHillSlopeRestorations()...)

	return modelActions
}

func (m *CoreModel) buildGullyRestorations() []action.ManagementAction {
	gullyRestorations := new(actions.GullyRestorationGroup).
		WithParameters(m.parameters).
		WithGullyTable(m.gulliesTable).
		WithActionsTable(m.actionsTable).
		ManagementActions()
	return gullyRestorations
}

func (m *CoreModel) buildRiverBankRestorations() []action.ManagementAction {
	riverBankRestorations := new(actions.RiverBankRestorationGroup).
		WithPlanningUnitTable(m.planningUnitTable).
		WithParentSoilsTable(m.actionsTable).
		WithParameters(m.parameters).
		ManagementActions()
	return riverBankRestorations
}

func (m *CoreModel) buildHillSlopeRestorations() []action.ManagementAction {
	hillSlopeRestorations := new(actions.HillSlopeRestorationGroup).
		WithPlanningUnitTable(m.planningUnitTable).
		WithActionsTable(m.actionsTable).
		WithParameters(m.parameters).
		ManagementActions()
	return hillSlopeRestorations
}

func (m *CoreModel) buildActionObservers() []action.Observer {
	observers := make([]action.Observer, 0)
	observers = append(observers, m)

	sedimentProduction2 := m.ContainedDecisionVariables.Variable(sedimentproduction.VariableName)
	if sedimentProduction2AsObserver, isObserver := sedimentProduction2.(action.Observer); isObserver {
		observers = append(observers, sedimentProduction2AsObserver)
	}

	nitrogenProduction := m.ContainedDecisionVariables.Variable(nitrogenproduction.VariableName)
	if nitrogenProductionAsObserver, isObserver := nitrogenProduction.(action.Observer); isObserver {
		observers = append(observers, nitrogenProductionAsObserver)
	}

	implementationCost := m.ContainedDecisionVariables.Variable(implementationcost.VariableName)
	if implementationCostAsObserver, isObserver := implementationCost.(action.Observer); isObserver {
		observers = append(observers, implementationCostAsObserver)
	}

	opportunityCost := m.ContainedDecisionVariables.Variable(opportunitycost.VariableName)
	if opportunityCostAsObserver, isObserver := opportunityCost.(action.Observer); isObserver {
		observers = append(observers, opportunityCostAsObserver)
	}

	return observers
}

func (m *CoreModel) observeActions(actionObservers []action.Observer, actions []action.ManagementAction) {
	m.managementActions.Initialise()
	for _, action := range actions {
		m.managementActions.Add(action)
		action.Subscribe(actionObservers...)
	}
}

func (m *CoreModel) randomlyInitialiseActions() {
	m.note("Starting randomly initialising model actions")

	m.initialising = true
	if m.parameters.HasEntry(parameters.MaximumImplementationCost) {
		m.note("Randomly initialising for Maximum implementation cost limit.")
		m.randomlyActivateActionsFromAllInactiveStart()
	} else if m.parameters.HasEntry(parameters.MaximumSedimentProduction) {
		m.note("Randomly initialising for Maximum sediment production limit.")
		m.randomlyDeactivateActionsFromAllActiveStart()
	} else if m.parameters.HasEntry(parameters.MaximumParticulateNitrogenProduction) {
		m.note("Randomly initialising for Maximum particulate nitrogen production limit.")
		m.randomlyDeactivateActionsFromAllActiveStart()
	} else {
		m.note("Randomly initialising for unbounded (no limits).")
		m.randomlyInitialiseActionsUnbounded()
	}
	m.initialising = false

	m.note("Finished randomly initialising model actions")
}

func (m *CoreModel) randomlyActivateActionsFromAllInactiveStart() {
	m.note("Initialising all actions as inactive")
	for _, action := range m.managementActions.Actions() {
		m.managementActions.RandomlyInitialiseAction(action)

		if !action.IsActive() {
			continue // Nothing to if it wasn't activated.
		}

		m.noteManagementAction("Randomly activating action", action)

		isValid, _ := m.ChangeIsValid()
		if !isValid {
			m.note("Change was invalid.  Reverting action to inactive")
			m.RevertChange()
			m.noteManagementAction("Action reverted", action)
		}
	}
}

func (m *CoreModel) randomlyDeactivateActionsFromAllActiveStart() {
	m.note("Initialising all actions as active")
	for _, action := range m.managementActions.Actions() {
		action.InitialisingActivation()
	}

	for _, action := range m.managementActions.Actions() {
		m.managementActions.RandomlyDeinitialiseAction(action)

		if action.IsActive() {
			continue // Nothing to if it wasn't deactivated.
		}

		m.noteManagementAction("Randomly deactivating action", action)

		isValid, _ := m.ChangeIsValid()
		if !isValid {
			m.note("Change was invalid.  Reverting to action being active")
			m.RevertChange()
			m.noteManagementAction("Action reverted", action)
		}
	}
}

func (m *CoreModel) randomlyInitialiseActionsUnbounded() {
	for _, action := range m.managementActions.Actions() {
		m.managementActions.RandomlyInitialiseAction(action)
		if action.IsActive() {
			m.noteManagementAction("Randomly activating action", action)
		}
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
	m.ContainedDecisionVariables.AcceptAll()
	if !m.initialising {
		m.noteManagementAction("Accepting Action", m.managementActions.LastAppliedAction())
	}
}

func (m *CoreModel) RevertChange() {
	if !m.initialising {
		m.noteManagementAction("Rejecting Action", m.managementActions.LastAppliedAction())
	}
	m.ContainedDecisionVariables.RejectAll()
	m.managementActions.ToggleLastActivationUnobserved()
}

func (m *CoreModel) DoRandomChange() {
	m.TryRandomChange()
	m.AcceptChange()
}

func (m *CoreModel) ToggleAction(planningUnit planningunit.Id, actionType action.ManagementActionType) {
	message := fmt.Sprintf("Toggling action [%v] for planning unit [%d]", actionType, planningUnit)
	m.note(message)
	m.managementActions.ToggleAction(planningUnit, actionType)
}

func (m *CoreModel) TryRandomChange() {
	m.managementActions.RandomlyToggleOneActivation()
	m.noteManagementAction("Trying Action", m.managementActions.LastAppliedAction())
}

func (m *CoreModel) UndoChange() {
	m.noteManagementAction("Undoing Action", m.managementActions.LastAppliedAction())
	m.managementActions.ToggleLastActivation()
}

func (m *CoreModel) ObserveAction(action action.ManagementAction) {
	// m.noteAppliedManagementAction(action)
}

func (m *CoreModel) ObserveActionInitialising(action action.ManagementAction) {
}

func (m *CoreModel) noteAppliedManagementAction(action action.ManagementAction) {
	if m.initialising {
		return
	}
	event := observer.NewEvent(observer.ManagementAction).
		WithAttribute("Type", action.Type()).
		WithAttribute("PlanningUnit", action.PlanningUnit()).
		WithAttribute("IsActive", action.IsActive())
	m.NotifyObserversOfEvent(*event)
}

func (m *CoreModel) note(text string) {
	event := observer.NewEvent(observer.Model).WithNote(text)
	m.NotifyObserversOfEvent(*event)
}

func (m *CoreModel) noteManagementAction(text string, action action.ManagementAction) {
	event := observer.NewEvent(observer.Model).
		WithNote(text).
		WithAttribute("Type", action.Type()).
		WithAttribute("PlanningUnit", action.PlanningUnit()).
		WithAttribute("IsActive", action.IsActive())
	m.NotifyObserversOfEvent(*event)
}

func (m *CoreModel) ObserveDecisionVariable(variable variable.DecisionVariable) {
	if m.initialising {
		return
	}
	event := observer.NewEvent(observer.DecisionVariable).
		WithAttribute("Name", variable.Name()).
		WithAttribute("Value", variable.Value())
	m.NotifyObserversOfEvent(*event)
}

func (m *CoreModel) ObserveDecisionVariableWithNote(variable variable.DecisionVariable, note string) {
	if m.initialising {
		return
	}
	event := observer.NewEvent(observer.DecisionVariable).
		WithAttribute("Name", variable.Name()).
		WithAttribute("Value", variable.Value()).
		WithNote(note)
	m.NotifyObserversOfEvent(*event)
}

func (m *CoreModel) capChangeOverRange(value float64) float64 {
	return math.Max(0, value)
}

func (m *CoreModel) DeepClone() model.Model {
	clone := *m
	clone.managementActions.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}

func (m *CoreModel) ChangeIsValid() (bool, *compositeErrors.CompositeError) {
	validationErrors := compositeErrors.New("Validation Errors")

	sedimentProduction := m.ContainedDecisionVariables.Variable(sedimentproduction.VariableName)
	if boundSedimentLoad, isBound := sedimentProduction.(variable.Bounded); isBound {
		if !boundSedimentLoad.WithinBounds(sedimentProduction.UndoableValue()) {
			validationMessage := fmt.Sprintf("SedimentProduction %s", boundSedimentLoad.BoundErrorAsText(sedimentProduction.UndoableValue()))
			validationErrors.AddMessage(validationMessage)
		}
	}

	implementationCost := m.ContainedDecisionVariables.Variable(implementationcost.VariableName)
	if boundImplementationCost, isBound := implementationCost.(variable.Bounded); isBound {
		if !boundImplementationCost.WithinBounds(implementationCost.UndoableValue()) {
			validationMessage := fmt.Sprintf("ImplementationCost value %s", boundImplementationCost.BoundErrorAsText(implementationCost.UndoableValue()))
			validationErrors.AddMessage(validationMessage)
		}
	}

	if validationErrors.Size() > 0 {
		return false, validationErrors
	}

	return true, nil
}

func (m *CoreModel) StateIsValid() (bool, *compositeErrors.CompositeError) {
	validationErrors := compositeErrors.New("Validation Errors")

	sedimentProduction := m.ContainedDecisionVariables.Variable(sedimentproduction.VariableName)
	if boundSedimentLoad, isBound := sedimentProduction.(variable.Bounded); isBound {
		if !boundSedimentLoad.WithinBounds(sedimentProduction.Value()) {
			validationMessage := fmt.Sprintf("SedimentProduction %s", boundSedimentLoad.BoundErrorAsText(sedimentProduction.Value()))
			validationErrors.AddMessage(validationMessage)
		}
	}

	implementationCost := m.ContainedDecisionVariables.Variable(implementationcost.VariableName)
	if boundImplementationCost, isBound := implementationCost.(variable.Bounded); isBound {
		if !boundImplementationCost.WithinBounds(implementationCost.Value()) {
			validationMessage := fmt.Sprintf("ImplementationCost value %s", boundImplementationCost.BoundErrorAsText(implementationCost.Value()))
			validationErrors.AddMessage(validationMessage)
		}
	}

	if validationErrors.Size() > 0 {
		return false, validationErrors
	}

	return true, nil
}

func (m *CoreModel) TearDown() {
	// deliberately does nothing.
}
