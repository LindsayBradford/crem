// Copyright (c) 2019 Australian Rivers Institute.

package action

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/pkg/errors"
	"sort"
)

// ModelManagementActions is a container/manager for all management actions that can be applied to a model.
type ModelManagementActions struct {
	lastApplied ManagementAction
	actions     ManagementActions
	rand.RandContainer
}

type ManagementActions []ManagementAction

func (m *ModelManagementActions) Initialise() {
	m.actions = make([]ManagementAction, 0)
	m.SetRandomNumberGenerator(rand.NewTimeSeeded())
}

// Add allows onr ore management actions to be added to the set of actions under management.
func (m *ModelManagementActions) Add(newActions ...ManagementAction) {
	for _, newAction := range newActions {
		m.actions = append(m.actions, newAction)
	}
}

func (m *ModelManagementActions) Sort() {
	sort.Sort(m.actions)
}

func (ma ManagementActions) Len() int {
	return len(ma)
}

func (ma ManagementActions) Swap(i, j int) {
	ma[i], ma[j] = ma[j], ma[i]
}

func (ma ManagementActions) Less(i, j int) bool {
	if ma[i].PlanningUnit() < ma[j].PlanningUnit() {
		return true
	}
	if ma[i].PlanningUnit() == ma[j].PlanningUnit() {
		if ma[i].Type() < ma[j].Type() {
			return true
		}
	}
	return false
}

// RandomlyToggleOneActivation randomly picks one of its stored management actions and toggles its activation
// in a way that will trigger any observers of the selected management action to react to its change in activation state.
func (m *ModelManagementActions) RandomlyToggleOneActivation() ManagementAction {
	m.lastApplied = m.pickRandomManagementAction()
	m.lastApplied.ToggleActivation()
	return m.lastApplied
}

func (m *ModelManagementActions) pickRandomManagementAction() ManagementAction {
	numberOfActions := len(m.actions)
	if numberOfActions < 1 {
		return NullManagementAction
	}
	randomIndex := m.RandomNumberGenerator().Intn(numberOfActions)
	return m.actions[randomIndex]
}

const (
	deactivate = 0
	activate   = 0
	ignore     = 1
)

// RandomlyInitialise will pass through its stored management actions, applying a 50/50 chance to activate
// each action. Any action chosen for activation triggers its observers to react to its 'initialising' activation.
func (m *ModelManagementActions) RandomlyInitialise() {
	for _, action := range m.actions {
		m.RandomlyInitialiseAction(action)
	}
}

func (m *ModelManagementActions) RandomlyInitialiseAction(action ManagementAction) {
	randomValue := m.RandomNumberGenerator().Intn(2)
	switch randomValue {
	case activate:
		m.lastApplied = action
		action.InitialisingActivation()
	case ignore:
		// Deliberately does nothing
	default:
		panic(errors.New("Random value outside range of [0,1]"))
	}
}

func (m *ModelManagementActions) RandomlyDeInitialiseAction(action ManagementAction) {
	randomValue := m.RandomNumberGenerator().Intn(2)
	switch randomValue {
	case deactivate:
		m.lastApplied = action
		action.InitialisingDeactivation()
	case ignore:
		// Deliberately does nothing
	default:
		panic(errors.New("Random value outside range of [0,1]"))
	}
}

func (m *ModelManagementActions) ToggleActionUnobserved(planningUnit planningunit.Id, actionType ManagementActionType) {
	for _, action := range m.actions {
		if action.PlanningUnit() == planningUnit && actionType == action.Type() {
			m.lastApplied = action
			m.ToggleLastActivationUnobserved()
		}
	}
}

func (m *ModelManagementActions) ToggleAction(planningUnit planningunit.Id, actionType ManagementActionType) {
	for _, action := range m.actions {
		if action.PlanningUnit() == planningUnit && actionType == action.Type() {
			m.lastApplied = action
			m.ToggleLastActivation()
		}
	}
}

// ToggleLastActivation allows for the last recorded management action change to have its
// activation state reverted, alerting any observers  of the change.
func (m *ModelManagementActions) ToggleLastActivation() {
	m.lastApplied.ToggleActivation()
}

// ToggleLastActivationUnobserved allows for the last recorded management action change to have its
// activation state reverted, without triggering any observation of the change.
func (m *ModelManagementActions) ToggleLastActivationUnobserved() {
	m.lastApplied.ToggleActivationUnobserved()
}

func (m *ModelManagementActions) DeactivateLastInitialisedAction() {
	m.lastApplied.InitialisingDeactivation()
}

func (m *ModelManagementActions) LastAppliedAction() ManagementAction {
	return m.lastApplied
}

func (m *ModelManagementActions) Actions() ManagementActions {
	return m.actions
}

func (m *ModelManagementActions) ActiveActions() ManagementActions {
	activeActions := make([]ManagementAction, 0)

	for _, action := range m.actions {
		if action.IsActive() {
			activeActions = append(activeActions, action)
		}
	}

	return activeActions
}

func (m *ModelManagementActions) SetActivation(index int, value bool) {
	m.actions[index].SetActivation(value)
}

func (m *ModelManagementActions) SetActivationUnobserved(index int, value bool) {
	m.actions[index].SetActivationUnobserved(value)
}
