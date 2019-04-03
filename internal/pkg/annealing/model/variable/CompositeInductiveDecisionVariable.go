// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"encoding/json"
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/action"
	errors2 "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/pkg/errors"
)

var _ InductiveDecisionVariable = new(CompositeInductiveDecisionVariable)

type CompositeInductiveDecisionVariable struct {
	name.NameContainer

	weightedVariables map[InductiveDecisionVariable]float64
	ContainedDecisionVariableObservers
	ContainedUnitOfMeasure
}

func (v *CompositeInductiveDecisionVariable) Initialise() *CompositeInductiveDecisionVariable {
	v.weightedVariables = make(map[InductiveDecisionVariable]float64, 0)
	return v
}

func (v *CompositeInductiveDecisionVariable) WithName(name string) *CompositeInductiveDecisionVariable {
	v.SetName(name)
	return v
}

func (v *CompositeInductiveDecisionVariable) WithWeightedVariable(variable InductiveDecisionVariable, weight float64) *CompositeInductiveDecisionVariable {
	v.weightedVariables[variable] = weight
	return v
}

func (v *CompositeInductiveDecisionVariable) Build() (*CompositeInductiveDecisionVariable, error) {
	errors := new(errors2.CompositeError)
	if weightError := v.checkWeights(); weightError != nil {
		errors.Add(weightError)
	}

	if errors.Size() > 0 {
		return nil, errors
	}

	return v, nil
}

func (v *CompositeInductiveDecisionVariable) vectorLengthOfVariableValues() float64 {
	var summedSquares float64
	for variable := range v.weightedVariables {
		summedSquares += math.Pow(variable.Value(), 2)
	}
	return math.Sqrt(summedSquares)
}

func (v *CompositeInductiveDecisionVariable) vectorLengthOfVariableInductiveValues() float64 {
	var summedSquares float64
	for variable := range v.weightedVariables {
		summedSquares += math.Pow(variable.InductiveValue(), 2)
	}
	return math.Sqrt(summedSquares)
}

func (v *CompositeInductiveDecisionVariable) checkWeights() error {
	overallWeight := float64(0)

	for _, weight := range v.weightedVariables {
		overallWeight += weight
	}

	if overallWeight == 1 {
		return nil
	}

	return errors.New("Variable weights do not add to one.")
}

func (v *CompositeInductiveDecisionVariable) Value() float64 {
	value := float64(0)
	for variable, weight := range v.weightedVariables {
		variableValue := variable.Value()
		// https://en.wikipedia.org/wiki/Feature_scaling#Scaling_to_unit_length
		scaledVariableValue := variableValue / v.vectorLengthOfVariableValues()
		weightedScaledValue := scaledVariableValue * weight
		value += weightedScaledValue
	}
	return value
}

func (v *CompositeInductiveDecisionVariable) SetValue(value float64) {
	// Deliberately does nothing
}

func (v *CompositeInductiveDecisionVariable) InductiveValue() float64 {
	value := float64(0)
	for variable, weight := range v.weightedVariables {
		variableValue := variable.InductiveValue()
		// https://en.wikipedia.org/wiki/Feature_scaling#Scaling_to_unit_length
		scaledVariableValue := variableValue / v.vectorLengthOfVariableInductiveValues()
		weightedScaledValue := scaledVariableValue * weight
		value += weightedScaledValue
	}
	return value
}

func (v *CompositeInductiveDecisionVariable) SetInductiveValue(value float64) {
	// Deliberately does nothing
}

func (v *CompositeInductiveDecisionVariable) AcceptInductiveValue() {
	for variable := range v.weightedVariables {
		variable.AcceptInductiveValue()
	}
	v.NotifyObservers()
}

func (v *CompositeInductiveDecisionVariable) RejectInductiveValue() {
	for variable := range v.weightedVariables {
		variable.RejectInductiveValue()
	}
	v.NotifyObservers()
}

func (v *CompositeInductiveDecisionVariable) DifferenceInValues() float64 {
	return v.InductiveValue() - v.Value()
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variable.
func (v *CompositeInductiveDecisionVariable) NotifyObservers() {
	for _, observer := range v.Observers() {
		observer.ObserveDecisionVariable(v)
	}
}

func (v *CompositeInductiveDecisionVariable) ObserveAction(action action.ManagementAction) {
	// deliberately does nothing
}

func (v *CompositeInductiveDecisionVariable) ObserveActionInitialising(action action.ManagementAction) {
	// deliberately does nothing
}

func (v *CompositeInductiveDecisionVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(MakeEncodeable(v))
}
