// Copyright (c) 2018 Australian Rivers Institute.

package dumb

import (
	"errors"
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
)

type Model struct {
	name.ContainedName
	rand.ContainedRand

	parameters        Parameters
	decisionVariables model.VolatileDecisionVariables
}

func New() *Model {
	newModel := new(Model)
	newModel.SetName("DumbModel")

	newModel.BuildDecisionVariables()
	newModel.parameters.Initialise()

	initialValue := newModel.parameters.GetFloat64(InitialObjectiveValue)
	newModel.decisionVariables[model.ObjectiveValue].SetValue(initialValue)

	return newModel
}

func (dm *Model) WithName(name string) *Model {
	dm.SetName(name)
	return dm
}

func (dm *Model) BuildDecisionVariables() {
	dm.decisionVariables = model.NewVolatileDecisionVariables()
	dm.createDecisionVariableFor(model.ObjectiveValue)
}

func (dm *Model) createDecisionVariableFor(name string) {
	newVariable := new(model.VolatileDecisionVariable)
	newVariable.SetName(name)
	dm.decisionVariables[name] = newVariable
}

func (dm *Model) WithParameters(params parameters.Map) *Model {
	dm.parameters.Merge(params)

	initialValue := dm.parameters.GetFloat64(InitialObjectiveValue)
	dm.decisionVariables[model.ObjectiveValue].SetValue(initialValue)

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
	return dm.decisionVariables[model.ObjectiveValue].Value()
}

func (dm *Model) setObjectiveValue(value float64) {
	dm.decisionVariables[model.ObjectiveValue].SetTemporaryValue(value)
}

func (dm *Model) SetDecisionVariable(name string, value float64) {
	dm.decisionVariables[name].SetValue(value)
}

func (dm *Model) AcceptChange() {
	dm.decisionVariables[model.ObjectiveValue].Accept()
}

func (dm *Model) RevertChange() {
	dm.decisionVariables[model.ObjectiveValue].Revert()
}

func (dm *Model) DecisionVariable(name string) (model.DecisionVariable, error) {
	if variable, found := dm.decisionVariables[name]; found == true {
		return variable, nil
	}
	return model.NullDecisionVariable, errors.New("decision variable [" + name + "] not defined for model [" + dm.Name() + " ].")
}

func (dm *Model) DecisionVariableChange(variableName string) (float64, error) {
	decisionVariable, foundActual := dm.decisionVariables[variableName]
	if !foundActual {
		return 0, errors.New("no temporary decision variable of name [" + decisionVariable.Name() + "] in model [" + dm.Name() + "].")
	}

	difference := decisionVariable.TemporaryValue() - decisionVariable.Value()
	return difference, nil
}

func (dm *Model) DeepClone() model.Model {
	clone := *dm
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
