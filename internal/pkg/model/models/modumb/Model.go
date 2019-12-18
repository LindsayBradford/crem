// Copyright (c) 2019 Australian Rivers Institute.

package modumb

import (
	"fmt"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb/variables"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	baseParameters "github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
)

type Model struct {
	name.NameContainer
	name.IdentifiableContainer
	rand.RandContainer

	parameters parameters.Parameters

	variable.ContainedDecisionVariables
	managementActions action.ModelManagementActions

	observer.ContainedEventNotifier
}

var Objectives = []string{
	"Objective_0",
	"Objective_1",
	"Objective_2",
}

func NewModel() *Model {
	newModel := new(Model)
	newModel.SetName("DumbMultiObjectiveModel")

	newModel.parameters.Initialise()
	newModel.ContainedDecisionVariables.Initialise()

	return newModel
}

func (m *Model) WithName(name string) *Model {
	m.SetName(name)
	return m
}

func (m *Model) WithId(id string) *Model {
	m.SetId(id)
	return m
}

func (m *Model) WithParameters(params baseParameters.Map) *Model {
	m.SetParameters(params)
	return m
}

func (m *Model) SetParameters(params baseParameters.Map) error {
	m.parameters.AssignAllUserValues(params)
	return m.ParameterErrors()
}

func (m *Model) ParameterErrors() error {
	return m.parameters.ValidationErrors()
}

func (m *Model) Initialise() {
	m.buildDecisionVariables()
	m.buildManagementActions()
}

func (m *Model) buildDecisionVariables() {
	m.ContainedDecisionVariables.Initialise()
	objectiveOne := new(variables.DumbObjective).
		Initialise().
		WithName(Objectives[0]).
		WithStartingValue(m.parameters.GetFloat64(parameters.InitialObjectiveOneValue)).
		WithObservers(m)

	objectiveTwo := new(variables.DumbObjective).
		Initialise().
		WithName(Objectives[1]).
		WithStartingValue(m.parameters.GetFloat64(parameters.InitialObjectiveTwoValue)).
		WithObservers(m)

	objectiveThree := new(variables.DumbObjective).
		Initialise().
		WithName(Objectives[2]).
		WithStartingValue(m.parameters.GetFloat64(parameters.InitialObjectiveThreeValue)).
		WithObservers(m)

	m.ContainedDecisionVariables.Add(
		objectiveOne, objectiveTwo, objectiveThree,
	)

	sortedKeys := m.ContainedDecisionVariables.DecisionVariables().SortedKeys()

	for _, value := range sortedKeys {
		variable := m.ContainedDecisionVariables.Variable(value)
		m.ObserveDecisionVariable(variable)
	}
}

func (m *Model) buildManagementActions() {
	m.managementActions.Initialise()
	numberOfPlanningUnits := m.parameters.GetInt64(parameters.NumberOfPlanningUnits)
	for planningUnit := planningunit.Id(0); planningUnit < planningunit.Id(numberOfPlanningUnits); planningUnit++ {

		actionOne := actions.New().
			WithObjectiveValue(action.ModelVariableName(Objectives[0]), -1).
			WithPlanningUnit(planningUnit)

		actionTwo := actions.New().
			WithObjectiveValue(action.ModelVariableName(Objectives[1]), -2).
			WithPlanningUnit(planningUnit)

		actionThree := actions.New().
			WithObjectiveValue(action.ModelVariableName(Objectives[2]), -3).
			WithPlanningUnit(planningUnit)

		m.managementActions.Add(actionOne, actionTwo, actionThree)

		objectiveOne := m.ContainedDecisionVariables.Variable(Objectives[0])
		if objectiveOneAsObserver, isObserver := objectiveOne.(action.Observer); isObserver {
			actionOne.Subscribe(m, objectiveOneAsObserver)
		}

		objectiveTwo := m.ContainedDecisionVariables.Variable(Objectives[1])
		if objectiveTwoAsObserver, isObserver := objectiveTwo.(action.Observer); isObserver {
			actionTwo.Subscribe(m, objectiveTwoAsObserver)
		}

		objectiveThree := m.ContainedDecisionVariables.Variable(Objectives[2])
		if objectiveThreeAsObserver, isObserver := objectiveThree.(action.Observer); isObserver {
			actionThree.Subscribe(m, objectiveThreeAsObserver)
		}
	}
}

func (m *Model) TearDown() {
	// This model doesn't need any special tearDown.
}

func (m *Model) DoRandomChange() {
	m.TryRandomChange()
	m.AcceptChange()
}

func (m *Model) UndoChange() {
	m.managementActions.ToggleLastActivation()
	m.AcceptChange()
}

func (m *Model) TryRandomChange() {
	m.note("Trying Random Change")
	m.managementActions.RandomlyToggleOneActivation()
}

func (m *Model) SetDecisionVariable(name string, value float64) {
	m.ContainedDecisionVariables.SetValue(name, value)
}

func (m *Model) AcceptChange() {
	m.ContainedDecisionVariables.AcceptAll()
}

func (m *Model) RevertChange() {
	m.ContainedDecisionVariables.RejectAll()
	m.managementActions.ToggleLastActivationUnobserved()
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

func (m *Model) SetManagementActionUnobserved(index int, value bool) {
	m.managementActions.SetActivationUnobserved(index, value)
}

func (m *Model) PlanningUnits() planningunit.Ids { return nil }

func (m *Model) DeepClone() model.Model {
	clone := *m
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
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

func (m *Model) ObserveAction(action action.ManagementAction) {
	m.noteAppliedManagementAction(action)
}

func (m *Model) ObserveActionInitialising(action action.ManagementAction) {
	m.noteAppliedManagementAction(action)
}

func (m *Model) noteAppliedManagementAction(actionToNote action.ManagementAction) {

	event := observer.NewEvent(observer.ManagementAction).
		WithId(m.Id()).
		WithAttribute("Type", actionToNote.Type()).
		WithAttribute("PlanningUnit", actionToNote.PlanningUnit()).
		WithAttribute("IsActive", actionToNote.IsActive())

	for _, name := range Objectives {
		modelVariableName := action.ModelVariableName(name)
		value := actionToNote.ModelVariableValue(modelVariableName)
		if value != 0 {
			noteText := fmt.Sprintf("Changing [%s] with active value=[%f]", name, value)
			event.WithNote(noteText)
		}
	}

	m.EventNotifier().NotifyObserversOfEvent(*event)
}
