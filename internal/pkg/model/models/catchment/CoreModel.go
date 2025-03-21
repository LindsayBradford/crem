// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	errors2 "errors"
	"fmt"
	catchmentDataSet "github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/dissolvednitrogen"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/opportunitycost"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/totalnitrogen"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/implementationcost"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/particulatenitrogen"
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

var _ model.Model = new(CoreModel)

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

	inputDataSet *catchmentDataSet.DataSetImpl
	initialising bool

	observer.SynchronousAnnealingEventNotifier

	attributes.ContainedAttributes

	managementActionEvent *observer.Event
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
	m.inputDataSet = new(catchmentDataSet.DataSetImpl).Initialise(sourceDataSet)
	return m
}

func (m *CoreModel) SetParameters(params baseParameters.Map) error {
	m.parameters.AssignAllUserValues(params)
	m.validateModelParameters()
	return m.parameters.ValidationErrors()
}

func (m *CoreModel) validateModelParameters() {
	boundVariableNumber := 0
	if m.parameters.HasEntry(parameters.MaximumSedimentProduction) {
		boundVariableNumber++
	}
	if m.parameters.HasEntry(parameters.MaximumParticulateNitrogenProduction) {
		boundVariableNumber++
	}
	if m.parameters.HasEntry(parameters.MaximumDissolvedNitrogenProduction) {
		boundVariableNumber++
	}
	if m.parameters.HasEntry(parameters.MaximumTotalNitrogenProduction) {
		boundVariableNumber++
	}
	if m.parameters.HasEntry(parameters.MaximumImplementationCost) {
		boundVariableNumber++
	}
	if m.parameters.HasEntry(parameters.MaximumOpportunityCost) {
		boundVariableNumber++
	}

	if boundVariableNumber > 1 {
		errorText := fmt.Sprintf("Only one of [%s], [%s], [%s], [%s], [%s] or [%s] allowed as variable limit.",
			parameters.MaximumSedimentProduction,
			parameters.MaximumParticulateNitrogenProduction,
			parameters.MaximumDissolvedNitrogenProduction,
			parameters.MaximumTotalNitrogenProduction,
			parameters.MaximumImplementationCost,
			parameters.MaximumOpportunityCost,
		)

		m.parameters.AddValidationErrorMessage(errorText)
	}
}

func (m *CoreModel) ParameterErrors() error {
	return m.parameters.ValidationErrors()
}

func (m *CoreModel) Initialise(initialisationType model.InitialisationType) {
	m.ReplaceAttribute("ModelSuppliedPlanningUnitName", "SubCatchment")
	m.planningUnitTable = m.fetchCsvTable(catchmentDataSet.SubcatchmentsTableName)
	m.gulliesTable = m.fetchCsvTable(catchmentDataSet.GulliesTableName)
	m.actionsTable = m.fetchCsvTable(catchmentDataSet.ActionsTableName)

	m.buildDecisionVariables()
	m.buildAndObserveManagementActions()
	m.InitialiseActions(initialisationType)

	firstAction := m.ManagementActions()[0]
	m.managementActionEvent = observer.NewEvent(observer.Model).
		WithNote("").
		WithAttribute("Type", firstAction.Type()).
		WithAttribute("PlanningUnit", firstAction.PlanningUnit()).
		WithAttribute("IsActive", firstAction.IsActive())
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
		Initialise(m.inputDataSet, m.parameters).
		WithObservers(m)

	if m.parameters.HasEntry(parameters.MaximumSedimentProduction) {
		sedimentProduction.SetMaximum(m.parameters.GetFloat64(parameters.MaximumSedimentProduction))
	}

	particulateNitrogen := new(particulatenitrogen.ParticulateNitrogenProduction).
		WithSedimentProductionVariable(sedimentProduction).
		Initialise(m.planningUnitTable, m.actionsTable, m.parameters).
		WithObservers(m)

	if m.parameters.HasEntry(parameters.MaximumParticulateNitrogenProduction) {
		particulateNitrogen.SetMaximum(m.parameters.GetFloat64(parameters.MaximumParticulateNitrogenProduction))
	}

	dissolvedNitrogen := new(dissolvednitrogen.DissolvedNitrogenProduction).
		Initialise(m.planningUnitTable, m.actionsTable, m.parameters).
		WithObservers(m)

	if m.parameters.HasEntry(parameters.MaximumDissolvedNitrogenProduction) {
		dissolvedNitrogen.SetMaximum(m.parameters.GetFloat64(parameters.MaximumDissolvedNitrogenProduction))
	}

	totalNitrogen := new(totalnitrogen.TotalNitrogenProduction).
		WithBaseNitrogenVariables(particulateNitrogen, dissolvedNitrogen).
		Initialise(m.planningUnitTable, m.actionsTable, m.parameters).
		WithObservers(m)

	if m.parameters.HasEntry(parameters.MaximumTotalNitrogenProduction) {
		totalNitrogen.SetMaximum(m.parameters.GetFloat64(parameters.MaximumTotalNitrogenProduction))
	}

	implementationCost := new(implementationcost.ImplementationCost).
		Initialise().WithObservers(m)

	if m.parameters.HasEntry(parameters.MaximumImplementationCost) {
		implementationCost.SetMaximum(m.parameters.GetFloat64(parameters.MaximumImplementationCost))
	}

	opportunityCost := new(opportunitycost.OpportunityCost).
		Initialise().WithObservers(m)

	if m.parameters.HasEntry(parameters.MaximumOpportunityCost) {
		opportunityCost.SetMaximum(m.parameters.GetFloat64(parameters.MaximumOpportunityCost))
	}

	m.ContainedDecisionVariables.Initialise()
	m.ContainedDecisionVariables.Add(
		sedimentProduction,
		particulateNitrogen, dissolvedNitrogen, totalNitrogen,
		implementationCost, opportunityCost,
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
	modelActions = append(modelActions, m.buildWetlandsEstablishments()...)

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
		WithActionsTable(m.actionsTable).
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

func (m *CoreModel) buildWetlandsEstablishments() []action.ManagementAction {
	wetlandsEstablishments := new(actions.WetlandsEstablishmentGroup).
		WithPlanningUnitTable(m.planningUnitTable).
		WithActionsTable(m.actionsTable).
		WithParameters(m.parameters).
		ManagementActions()
	return wetlandsEstablishments
}

func (m *CoreModel) buildActionObservers() []action.Observer {
	observers := make([]action.Observer, 0)
	observers = append(observers, m)

	// Order of variable creation is crucial here for having downstream dependent observers
	// trigger once upstream variables have processed an action.
	for _, variable := range m.ContainedDecisionVariables.CreationOrderedVariables() {
		if variableAsObserver, isObserver := variable.(action.Observer); isObserver {
			observers = append(observers, variableAsObserver)
		}
	}

	return observers
}

func (m *CoreModel) observeActions(actionObservers []action.Observer, actions []action.ManagementAction) {
	m.managementActions.Initialise()
	for _, action := range actions {
		m.managementActions.Add(action)
		action.Subscribe(actionObservers...)
	}
	m.managementActions.Sort()
}

func (m *CoreModel) InitialiseActions(initialisationType model.InitialisationType) {
	m.note("Starting initialising model actions")

	m.initialising = true

	if initialisationType == model.Unchanged {
		m.note("Initialising to an unchanged state.")
	} else if initialisationType == model.AsIs {
		m.note("Initialising to As-Is state.")
		m.InitialiseAllActionsToInactive()
	} else if m.parameters.HasEntry(parameters.MaximumImplementationCost) {
		m.note("Initialising for Maximum implementation cost limit.")
		m.InitialiseAllActionsToInactive()
	} else if m.parameters.HasEntry(parameters.MaximumOpportunityCost) {
		m.note("Initialising for Maximum opportunity cost limit.")
		m.InitialiseAllActionsToInactive()
	} else if m.parameters.HasEntry(parameters.MaximumSedimentProduction) {
		m.note("Randomly initialising for Maximum sediment production limit.")
		m.InitialiseAllActionsToActive()
	} else if m.parameters.HasEntry(parameters.MaximumParticulateNitrogenProduction) {
		m.note("Randomly initialising for Maximum particulate nitrogen production limit.")
		m.InitialiseAllActionsToActive()
	} else if m.parameters.HasEntry(parameters.MaximumDissolvedNitrogenProduction) {
		m.note("Randomly initialising for Maximum dissolved nitrogen production limit.")
		m.InitialiseAllActionsToActive()
	} else if m.parameters.HasEntry(parameters.MaximumTotalNitrogenProduction) {
		m.note("Randomly initialising for Maximum total nitrogen production limit.")
		m.InitialiseAllActionsToActive()
	}

	m.initialising = false

	m.note("Finished initialising model actions")
}

func (m *CoreModel) Randomize() {
	m.note("Starting randomizing model action state")

	if m.parameters.HasEntry(parameters.MaximumImplementationCost) {
		m.note("Randomly initialising for Maximum implementation cost limit.")
		m.RandomlyValidlyActivateActions()
	} else if m.parameters.HasEntry(parameters.MaximumOpportunityCost) {
		m.note("Randomly initialising for Maximum opportunity cost limit.")
		m.RandomlyValidlyActivateActions()
	} else if m.parameters.HasEntry(parameters.MaximumSedimentProduction) {
		m.note("Randomly initialising for Maximum sediment production limit.")
		m.RandomlyValidlyDeactivateActions()
	} else if m.parameters.HasEntry(parameters.MaximumParticulateNitrogenProduction) {
		m.note("Randomly initialising for Maximum particulate nitrogen production limit.")
		m.RandomlyValidlyDeactivateActions()
	} else if m.parameters.HasEntry(parameters.MaximumDissolvedNitrogenProduction) {
		m.note("Randomly initialising for Maximum dissolved nitrogen production limit.")
		m.RandomlyValidlyDeactivateActions()
	} else if m.parameters.HasEntry(parameters.MaximumTotalNitrogenProduction) {
		m.note("Randomly initialising for Maximum total nitrogen production limit.")
		m.RandomlyValidlyDeactivateActions()
	} else {
		m.note("Randomly initialising for unbounded (no limits).")
		m.randomlyInitialiseActionsUnbounded()
	}
	m.note("Finished randomizing model action state")
}

func (m *CoreModel) randomlyActivateActionsFromAllInactiveStart() {
	m.InitialiseAllActionsToInactive()
	m.RandomlyValidlyActivateActions()
}

func (m *CoreModel) InitialiseAllActionsToInactive() {
	m.note("Initialising all actions as inactive")
	for _, action := range m.managementActions.Actions() {
		action.InitialisingDeactivation()
	}
}

func (m *CoreModel) RandomlyValidlyActivateActions() {
	actionNumber := len(m.managementActions.Actions())
	attemptLimit := actionNumber

	attemptNote := fmt.Sprintf("Making [%d] attempts to find solution close to decision variable limit.", attemptLimit)
	m.note(attemptNote)

	isValid := true
	for isValid && attemptLimit > 0 {
		actionChanged := m.managementActions.RandomlyInitialiseAnyAction()
		if actionChanged == nil {
			continue
		}

		m.noteManagementAction("Activated random action", actionChanged)
		isValid, _ = m.ChangeIsValid()

		if isValid {
			attemptNote = fmt.Sprintf("Attempt [%d]: Action was valid. Keeping.", actionNumber-attemptLimit+1)
			m.note(attemptNote)
		} else {
			attemptNote = fmt.Sprintf("Attempt [%d]: Activation was invalid. Reverting.", actionNumber-attemptLimit+1)
			m.note(attemptNote)
			actionChanged.InitialisingDeactivation()
		}
		attemptLimit--
	}
	if attemptLimit == 0 {
		attemptNote := "Attempt limit reached while seeking a solution near configured decision variable limit. Please check configuration."
		m.note(attemptNote)

		panic(errors2.New(attemptNote))
	} else {
		m.note("Solution close to limit found. Using this as initial model state.")
	}
}

func (m *CoreModel) randomlyDeactivateActionsFromAllActiveStart() {
	m.InitialiseAllActionsToActive()
	m.RandomlyValidlyDeactivateActions()
}

func (m *CoreModel) InitialiseAllActionsToActive() {
	m.note("Initialising all actions as active")
	for _, action := range m.managementActions.Actions() {
		action.InitialisingActivation()
	}
}

func (m *CoreModel) RandomlyValidlyDeactivateActions() {
	numberToAttempt := len(m.managementActions.Actions())
	attemptsLeft := numberToAttempt

	attemptNote := fmt.Sprintf("Making [%d] attempts to find solution close to decision variable limit.", attemptsLeft)
	m.note(attemptNote)

	isValid := true
	for isValid && attemptsLeft > 0 {
		actionChanged := m.managementActions.RandomlyDeInitialiseAnyAction()
		if actionChanged == nil {
			continue
		}

		m.noteManagementAction("Deactivate random action", actionChanged)
		isValid, _ = m.ChangeIsValid()

		if isValid {
			attemptNote = fmt.Sprintf("Attempt [%d]: Action was valid. Keeping.", numberToAttempt-attemptsLeft+1)
			m.note(attemptNote)
		} else {
			attemptNote = fmt.Sprintf("Attempt [%d]: Activation was invalid. Reverting.", numberToAttempt-attemptsLeft+1)
			m.note(attemptNote)
			actionChanged.InitialisingActivation()
		}
		attemptsLeft--
	}
	if attemptsLeft == 0 {
		attemptNote := "Attempt limit reached while seeking a solution near configured decision variable limit. Please check configuration."
		m.note(attemptNote)

		panic(errors2.New(attemptNote))
	} else {
		m.note("Solution close to limit found. Using this as initial solution.")
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
	if m.ManagementActions()[index].IsActive() != value {
		m.managementActions.SetActivation(index, value)
		m.AcceptChange()
	}
}

func (m *CoreModel) SetManagementActionUnobserved(index int, value bool) {
	if m.ManagementActions()[index].IsActive() != value {
		m.managementActions.SetActivationUnobserved(index, value)
		m.AcceptChange()
	}
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

	if !m.HasObservers() {
		return
	}
	event := observer.NewEvent(observer.ManagementAction).
		WithAttribute("Type", action.Type()).
		WithAttribute("PlanningUnit", action.PlanningUnit()).
		WithAttribute("IsActive", action.IsActive())
	m.NotifyObserversOfEvent(*event)
}

func (m *CoreModel) note(text string) {
	if !m.HasObservers() {
		return
	}

	event := observer.NewEvent(observer.Model).WithNote(text)
	m.NotifyObserversOfEvent(*event)
}

func (m *CoreModel) noteManagementAction(text string, action action.ManagementAction) {
	if !m.HasObservers() {
		return
	}

	m.managementActionEvent.
		ReplaceNote(text).
		ReplacingAttribute("Type", action.Type()).
		ReplacingAttribute("PlanningUnit", action.PlanningUnit()).
		ReplaceAttribute("IsActive", action.IsActive())

	m.NotifyObserversOfEvent(*m.managementActionEvent)
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
	return m.checkValidityWith(m.undoableValueBoundsChecker)
}

func (m *CoreModel) checkValidityWith(validationFunction func(*compositeErrors.CompositeError)) (bool, *compositeErrors.CompositeError) {
	validationErrors := compositeErrors.New("Validation Errors")

	validationFunction(validationErrors)
	if validationErrors.Size() > 0 {
		return false, validationErrors
	}

	return true, nil
}

func (m *CoreModel) undoableValueBoundsChecker(validationErrors *compositeErrors.CompositeError) {
	variableNames := m.DecisionVariableNames()
	for _, name := range variableNames {
		m.checkVariableUndoableValueBounds(name, validationErrors)
	}
}

func (m *CoreModel) checkVariableUndoableValueBounds(variableName string, validationErrors *compositeErrors.CompositeError) {
	variableToCheck := m.ContainedDecisionVariables.Variable(variableName)
	checkBounds(variableToCheck, variableToCheck.UndoableValue(), validationErrors)
}

func (m *CoreModel) StateIsValid() (bool, *compositeErrors.CompositeError) {
	return m.checkValidityWith(m.actualValueBoundsChecker)
}

func (m *CoreModel) actualValueBoundsChecker(validationErrors *compositeErrors.CompositeError) {
	variableNames := m.DecisionVariableNames()
	for _, name := range variableNames {
		m.checkVariableActualValueBounds(name, validationErrors)
	}
}

func (m *CoreModel) checkVariableActualValueBounds(variableName string, validationErrors *compositeErrors.CompositeError) {
	variableToCheck := m.ContainedDecisionVariables.Variable(variableName)
	checkBounds(variableToCheck, variableToCheck.Value(), validationErrors)
}

func checkBounds(possiblyBoundVariable variable.UndoableDecisionVariable, value float64, validationErrors *compositeErrors.CompositeError) {
	if boundVariable, isBound := possiblyBoundVariable.(variable.Bounded); isBound {
		if !boundVariable.WithinBounds(value) {
			message := fmt.Sprintf("%s %s", possiblyBoundVariable.Name(), boundVariable.BoundErrorAsText(value))
			validationErrors.AddMessage(message)
		}
	}
}

func (m *CoreModel) TearDown() {
	// deliberately does nothing.
}

func (m *CoreModel) IsEquivalentTo(otherModel model.Model) bool {
	if !m.checkActions(otherModel) {
		return false
	}
	if !m.checkVariables(otherModel) {
		return false
	}
	return true
}

func (m *CoreModel) checkActions(otherModel model.Model) bool {
	myActions := m.ManagementActions()
	otherActions := otherModel.ManagementActions()
	for index := range myActions {
		assert.That(myActions[index].PlanningUnit() == otherActions[index].PlanningUnit()).Holds()
		assert.That(myActions[index].Type() == otherActions[index].Type()).Holds()

		if myActions[index].IsActive() != otherActions[index].IsActive() {
			return false
		}
	}
	return true
}

func (m *CoreModel) checkVariables(otherModel model.Model) bool {
	myDecisionVariables := *m.NameMappedVariables()
	for _, variable := range myDecisionVariables {
		otherVariable := otherModel.DecisionVariable(variable.Name())
		if variable.Value() != otherVariable.Value() {
			return false
		}
	}

	return true
}

func (m *CoreModel) SynchroniseTo(otherModel model.Model) {
	for index, action := range otherModel.ManagementActions() {
		m.SetManagementAction(index, action.IsActive())
	}
}
