// Copyright (c) 2018 Australian Rivers Institute.

package dumb

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
)

type Model struct {
	name.ContainedName
	rand.ContainedRand

	parameters        Parameters
	decisionVariables DecisionVariables
}

func New() *Model {
	newModel := new(Model)
	newModel.SetName("DumbModel")

	newModel.decisionVariables.Initialise()
	newModel.parameters.Initialise()

	initialValue := newModel.parameters.GetFloat64(InitialObjectiveValue)
	newModel.decisionVariables.SetValue(model.ObjectiveValue, initialValue)

	return newModel
}

func (dm *Model) WithName(name string) *Model {
	dm.SetName(name)
	return dm
}

func (dm *Model) WithParameters(params parameters.Map) *Model {
	dm.parameters.Merge(params)

	initialValue := dm.parameters.GetFloat64(InitialObjectiveValue)
	dm.decisionVariables.SetValue(model.ObjectiveValue, initialValue)

	return dm
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
	return dm.decisionVariables.Value(model.ObjectiveValue)
}

func (dm *Model) setObjectiveValue(value float64) {
	dm.decisionVariables.Variable(model.ObjectiveValue).SetTemporaryValue(value)
}

func (dm *Model) SetDecisionVariable(name string, value float64) {
	dm.decisionVariables.SetValue(name, value)
}

func (dm *Model) AcceptChange() {
	dm.decisionVariables.Variable(model.ObjectiveValue).Accept()
}

func (dm *Model) RevertChange() {
	dm.decisionVariables.Variable(model.ObjectiveValue).Revert()
}

func (dm *Model) DecisionVariable(name string) model.DecisionVariable {
	return dm.decisionVariables.Variable(name)
}

func (dm *Model) DecisionVariableChange(variableName string) float64 {
	decisionVariable := dm.decisionVariables.Variable(variableName)
	difference := decisionVariable.TemporaryValue() - decisionVariable.Value()
	return difference
}

func (dm *Model) DeepClone() model.Model {
	clone := *dm
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
