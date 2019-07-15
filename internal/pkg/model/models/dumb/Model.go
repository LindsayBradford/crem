// Copyright (c) 2018 Australian Rivers Institute.

package dumb

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
)

type Model struct {
	name.NameContainer
	name.IdentifiableContainer
	rand.RandContainer

	parameters Parameters
	variable.ContainedDecisionVariables
}

func NewModel() *Model {
	newModel := new(Model)
	newModel.SetName("DumbModel")

	newModel.DecisionVariables()
	newModel.parameters.Initialise()

	newModel.ContainedDecisionVariables.Initialise()
	newModel.ContainedDecisionVariables.NewForName(variable.ObjectiveValue)

	initialValue := newModel.parameters.GetFloat64(InitialObjectiveValue)
	newModel.ContainedDecisionVariables.SetValue(variable.ObjectiveValue, initialValue)

	return newModel
}

func (dm *Model) WithName(name string) *Model {
	dm.SetName(name)
	return dm
}

func (dm *Model) WithParameters(params parameters.Map) *Model {
	dm.SetParameters(params)

	return dm
}

func (dm *Model) SetParameters(params parameters.Map) error {
	dm.parameters.AssignAllUserValues(params)

	initialValue := dm.parameters.GetFloat64(InitialObjectiveValue)
	dm.ContainedDecisionVariables.SetValue(variable.ObjectiveValue, initialValue)

	return dm.ParameterErrors()
}

func (dm *Model) ParameterErrors() error {
	return dm.parameters.ValidationErrors()
}

const (
	downward = -1
	upward   = 1
)

func (dm *Model) Initialise() {
	dm.SetRandomNumberGenerator(rand.NewTimeSeeded())
}

func (dm *Model) TearDown() {
	// This model doesn't need any special tearDown.
}

func (dm *Model) TryRandomChange() {
	originalValue := dm.objectiveValue()
	change := dm.generateRandomChange()
	newValue := dm.capChangeOverRange(originalValue + change)
	dm.setObjectiveValue(newValue)
}

func (dm *Model) generateRandomChange() float64 {
	randomValue := dm.RandomNumberGenerator().Intn(2)

	var changeInObjectiveValue float64
	switch randomValue {
	case 0:
		changeInObjectiveValue = downward
	case 1:
		changeInObjectiveValue = upward
	}

	return changeInObjectiveValue
}

func (dm *Model) capChangeOverRange(value float64) float64 {
	maxCappedValue := math.Max(dm.parameters.GetFloat64(MinimumObjectiveValue), value)
	bothCappedValue := math.Min(dm.parameters.GetFloat64(MaximumObjectiveValue), maxCappedValue)
	return bothCappedValue
}

func (dm *Model) objectiveValue() float64 {
	return dm.ContainedDecisionVariables.Value(variable.ObjectiveValue)
}

func (dm *Model) setObjectiveValue(value float64) {
	dm.ContainedDecisionVariables.Variable(variable.ObjectiveValue).SetInductiveValue(value)
}

func (dm *Model) SetDecisionVariable(name string, value float64) {
	dm.ContainedDecisionVariables.SetValue(name, value)
}

func (dm *Model) AcceptChange() {
	dm.ContainedDecisionVariables.Variable(variable.ObjectiveValue).AcceptInductiveValue()
}

func (dm *Model) RevertChange() {
	dm.ContainedDecisionVariables.Variable(variable.ObjectiveValue).RejectInductiveValue()
}

func (dm *Model) ManagementActions() []action.ManagementAction        { return nil }
func (dm *Model) ActiveManagementActions() []action.ManagementAction  { return nil }
func (dm *Model) SetManagementAction(index int, value bool)           {}
func (dm *Model) SetManagementActionUnobserved(index int, value bool) {}

func (dm *Model) PlanningUnits() planningunit.Ids { return nil }

func (dm *Model) DeepClone() model.Model {
	clone := *dm
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
