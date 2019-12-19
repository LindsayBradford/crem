// Copyright (c) 2018 Australian Rivers Institute.

package dumb

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/pkg/errors"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
)

type Model struct {
	name.NameContainer
	name.IdentifiableContainer
	rand.RandContainer

	parameters Parameters
	variable   *variable.SimpleUndoableDecisionVariable
}

func NewModel() *Model {
	newModel := new(Model)
	newModel.SetName("DumbModel")

	newModel.parameters.Initialise()
	newModel.variable = variable.NewUndoableDecisionVariable("ObjectiveValue")
	initialValue := newModel.parameters.GetFloat64(InitialObjectiveValue)
	newModel.variable.SetValue(initialValue)

	return newModel
}

func (m *Model) WithName(name string) *Model {
	m.SetName(name)
	return m
}

func (m *Model) WithParameters(params parameters.Map) *Model {
	m.SetParameters(params)

	return m
}

func (m *Model) SetParameters(params parameters.Map) error {
	m.parameters.AssignAllUserValues(params)

	initialValue := m.parameters.GetFloat64(InitialObjectiveValue)
	m.variable.SetValue(initialValue)

	return m.ParameterErrors()
}

func (m *Model) ParameterErrors() error {
	return m.parameters.ValidationErrors()
}

const (
	downward = -1
	upward   = 1
)

func (m *Model) Initialise() {
	m.SetRandomNumberGenerator(rand.NewTimeSeeded())
}

func (m *Model) TearDown() {
	// This model doesn't need any special tearDown.
}

func (m *Model) DoRandomChange() {
	m.TryRandomChange()
	m.AcceptChange()
}

func (m *Model) UndoChange() {
	m.variable.ApplyUndoneValue()
}

func (m *Model) TryRandomChange() {
	change := m.generateRandomChange()
	m.variable.SetUndoableChange(change)
}

func (m *Model) generateRandomChange() float64 {
	randomValue := m.RandomNumberGenerator().Intn(2)

	var changeInObjectiveValue float64
	switch randomValue {
	case 0:
		changeInObjectiveValue = downward
	case 1:
		changeInObjectiveValue = upward
	}

	return changeInObjectiveValue
}

func (m *Model) capChangeOverRange(value float64) float64 {
	maxCappedValue := math.Max(m.parameters.GetFloat64(MinimumObjectiveValue), value)
	bothCappedValue := math.Min(m.parameters.GetFloat64(MaximumObjectiveValue), maxCappedValue)
	return bothCappedValue
}

func (m *Model) objectiveValue() float64 {
	return m.variable.Value()
}

func (m *Model) SetDecisionVariable(name string, value float64) {
	m.variable.SetValue(value)
}

func (m *Model) AcceptChange() {
	m.variable.ApplyDoneValue()
}

func (m *Model) RevertChange() {
	m.variable.ApplyUndoneValue()
}

func (m *Model) ManagementActions() []action.ManagementAction        { return nil }
func (m *Model) ActiveManagementActions() []action.ManagementAction  { return nil }
func (m *Model) SetManagementAction(index int, value bool)           {}
func (m *Model) SetManagementActionUnobserved(index int, value bool) {}

func (m *Model) PlanningUnits() planningunit.Ids { return nil }

func (m *Model) DeepClone() model.Model {
	clone := *m
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}

func (m *Model) DecisionVariables() *variable.DecisionVariableMap {
	varMap := make(variable.DecisionVariableMap, 1)
	varMap["ObjectiveValue"] = m.variable
	return &varMap
}

func (m *Model) DecisionVariable(name string) variable.DecisionVariable {
	return m.variable
}

func (m *Model) OffersDecisionVariable(name string) bool {
	if name == "ObjectiveValue" {
		return true
	}
	return false
}

func (m *Model) DecisionVariableChange(decisionVariableName string) float64 {
	return m.variable.DifferenceInValues()
}

func (m *Model) ChangeIsValid() (bool, *errors.CompositeError) { return true, nil }
