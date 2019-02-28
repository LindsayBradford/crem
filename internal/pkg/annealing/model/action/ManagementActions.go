// Copyright (c) 2019 Australian Rivers Institute.

package action

import (
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/pkg/errors"
)

type ManagementActions struct {
	lastApplied ManagementAction
	actions     []ManagementAction
	rand.ContainedRand
}

func (m *ManagementActions) Initialise() {
	m.actions = make([]ManagementAction, 0)
	m.SetRandomNumberGenerator(rand.NewTimeSeeded())
}

func (m *ManagementActions) Add(newAction ManagementAction) {
	m.actions = append(m.actions, newAction)
}

func (m *ManagementActions) RandomlyToggleOneActivation() {
	m.lastApplied = m.pickRandomManagementAction()
	m.lastApplied.ToggleActivation()
}

func (m *ManagementActions) pickRandomManagementAction() ManagementAction {
	numberOfActions := len(m.actions)
	if numberOfActions < 1 {
		return NullManagementAction
	}
	randomIndex := m.RandomNumberGenerator().Intn(numberOfActions)
	return m.actions[randomIndex]
}

func (m *ManagementActions) RandomlyToggleAllActivations() {
	for _, action := range m.actions {
		randomValue := m.RandomNumberGenerator().Intn(2)
		switch randomValue {
		case 0:
			action.ToggleInitialisingActivation()
		case 1:
			// Deliberately does nothing
		default:
			panic(errors.New("Random value outside range of [0,1]"))
		}
	}
}

func (m *ManagementActions) UndoLastActivationToggleUnobserved() {
	m.lastApplied.ToggleActivationUnobserved()
}
