// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	errors2 "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"

	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	baseParameters "github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/LindsayBradford/crem/pkg/strings"
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

	m.buildCoreDecisionVariables()
	m.buildManagementActions()

	m.randomlyInitialiseActions()

	m.buildSedimentVsCostDecisionVariable()
}

func (m *CoreModel) randomlyInitialiseActions() {
	validCombinationFound := false
	for _, action := range m.managementActions.Actions() {
		m.managementActions.RandomlyInitialiseAction(action)
		isValid, _ := m.ChangeIsValid()

		if !validCombinationFound && isValid {
			validCombinationFound = true
			m.note("Found at least one valid scenario for specified variable limits")
			continue
		}

		if validCombinationFound && !isValid {
			m.note("Scenario would be invalid, reverting to last valid solution")
			m.managementActions.DeactivateLastInitialisedAction()
		}
	}
}

func (m *CoreModel) fetchPlanningUnitTable() tables.CsvTable {
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

func (m *CoreModel) fetchGulliesTable() tables.CsvTable {
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

func (m *CoreModel) buildCoreDecisionVariables() {
	sedimentLoad := new(variables.SedimentProduction).
		Initialise(m.planningUnitTable, m.gulliesTable, m.parameters).
		WithObservers(m)

	if m.parameters.HasEntry(parameters.MinimumSedimentProduction) {
		sedimentLoad.SetMinimum(m.parameters.GetFloat64(parameters.MinimumSedimentProduction))
	}

	if m.parameters.HasEntry(parameters.MaximumSedimentProduction) {
		sedimentLoad.SetMaximum(m.parameters.GetFloat64(parameters.MaximumSedimentProduction))
	}

	implementationCost := new(variables.ImplementationCost).
		Initialise(m.planningUnitTable, m.parameters).
		WithObservers(m)

	if m.parameters.HasEntry(parameters.MinimumImplementationCost) {
		implementationCost.SetMinimum(m.parameters.GetFloat64(parameters.MinimumImplementationCost))
	}

	if m.parameters.HasEntry(parameters.MaximumImplementationCost) {
		implementationCost.SetMaximum(m.parameters.GetFloat64(parameters.MaximumImplementationCost))
	}

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

func (m *CoreModel) buildSedimentVsCostDecisionVariable() {
	sedimentProduction := m.ContainedDecisionVariables.Variable(variables.SedimentProductionVariableName)
	implementationCost := m.ContainedDecisionVariables.Variable(variables.ImplementationCostVariableName)

	sedimentWeight := m.parameters.GetFloat64(parameters.SedimentProductionDecisionWeight)
	implementationCostWeight := m.parameters.GetFloat64(parameters.ImplementationCostDecisionWeight)

	sedimentVsCost, buildError := new(variables.SedimentVsCost).
		Initialise().
		WithObservers(m).
		WithWeightedVariable(sedimentProduction, sedimentWeight).
		WithWeightedVariable(implementationCost, implementationCostWeight).
		Build()

	if buildError != nil {
		panic(buildError)
	}

	noteBuilder := new(strings.FluentBuilder).
		Add(sedimentProduction.Name(), " weight = ", strconv.FormatFloat(sedimentWeight, 'f', 3, 64), ", ").
		Add(implementationCost.Name(), " weight = ", strconv.FormatFloat(implementationCostWeight, 'f', 3, 64))

	m.ObserveDecisionVariableWithNote(sedimentProduction, " Initial Value")
	m.ObserveDecisionVariableWithNote(implementationCost, " Initial Value")
	m.ObserveDecisionVariableWithNote(sedimentVsCost, noteBuilder.String())

	m.ContainedDecisionVariables.Add(sedimentVsCost)
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

func (m *CoreModel) ChangeIsValid() (bool, *errors2.CompositeError) {
	validationErrors := errors2.New("Validation Errors")

	sedimentProduction := m.ContainedDecisionVariables.Variable(variables.SedimentProductionVariableName)
	if boundSedimentLoad, isSedimentBound := sedimentProduction.(variable.Bounded); isSedimentBound {
		if !boundSedimentLoad.WithinBounds(sedimentProduction.InductiveValue()) {
			validationMessage := fmt.Sprintf("SedimentProduction %s", boundSedimentLoad.BoundErrorAsText(sedimentProduction.InductiveValue()))
			validationErrors.AddMessage(validationMessage)
		}
	}

	implementationCost := m.ContainedDecisionVariables.Variable(variables.ImplementationCostVariableName)
	if boundImplementationCost, isCostBound := implementationCost.(variable.Bounded); isCostBound {
		if !boundImplementationCost.WithinBounds(implementationCost.InductiveValue()) {
			validationMessage := fmt.Sprintf("ImplementationCost value %s", boundImplementationCost.BoundErrorAsText(implementationCost.InductiveValue()))
			validationErrors.AddMessage(validationMessage)
		}
	}

	if validationErrors.Size() > 0 {
		return false, validationErrors
	}

	return true, nil
}

func (m *CoreModel) deriveDataSourcePath() string {
	relativeFilePath := m.parameters.GetString(parameters.DataSourcePath)
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, relativeFilePath)
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
		WithAttribute("PlanningUnit", action.PlanningUnit()).
		WithAttribute("IsActive", action.IsActive())
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

func (m *CoreModel) TearDown() {
	// Deliberately does nothing
}
