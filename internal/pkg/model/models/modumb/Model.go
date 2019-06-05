// Copyright (c) 2019 Australian Rivers Institute.

package modumb

import (
	"fmt"

	baseParameters "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb/variables"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
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
	m.parameters.Merge(params)
	return m
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
		WithName("ObjectiveOne").
		WithStartingValue(m.parameters.GetFloat64(parameters.InitialObjectiveOneValue)).
		WithObservers(m)

	objectiveTwo := new(variables.DumbObjective).
		Initialise().
		WithName("ObjectiveTwo").
		WithStartingValue(m.parameters.GetFloat64(parameters.InitialObjectiveTwoValue)).
		WithObservers(m)

	objectiveThree := new(variables.DumbObjective).
		Initialise().
		WithName("ObjectiveThree").
		WithStartingValue(m.parameters.GetFloat64(parameters.InitialObjectiveThreeValue)).
		WithObservers(m)

	m.ContainedDecisionVariables.Add(
		objectiveOne, objectiveTwo, objectiveThree,
	)
}

func (m *Model) buildManagementActions() {
	m.managementActions.Initialise()
	numberOfPlanningUnits := m.parameters.GetInt64(parameters.NumberOfPlanningUnits)
	for planningUnit := int64(0); planningUnit < numberOfPlanningUnits; planningUnit++ {

		planningUnitAsString := fmt.Sprintf("%d", planningUnit)

		actionOne := actions.New().
			WithObjectiveValue("ObjectiveOne", 1).
			WithPlanningUnit(planningUnitAsString)

		actionTwo := actions.New().
			WithObjectiveValue("ObjectiveOne", 2).
			WithPlanningUnit(planningUnitAsString)

		actionThree := actions.New().
			WithObjectiveValue("ObjectiveOne", 4).
			WithPlanningUnit(planningUnitAsString)

		m.managementActions.Add(actionOne, actionTwo, actionThree)

		actionOne.Subscribe(m.ContainedDecisionVariables.Variable("ObjectiveOne"))
		actionTwo.Subscribe(m.ContainedDecisionVariables.Variable("ObjectiveTwo"))
		actionThree.Subscribe(m.ContainedDecisionVariables.Variable("ObjectiveThree"))
	}
}

func (m *Model) TearDown() {
	// This model doesn't need any special tearDown.
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
	m.managementActions.UndoLastActivationToggleUnobserved()
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

func (m *Model) PlanningUnits() solution.PlanningUnitIds { return nil }

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
