// Copyright (c) 2019 Australian Rivers Institute.

package variableOld

import (
	"encoding/json"
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	errors2 "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/pkg/errors"
)

var _ variable.UndoableDecisionVariable = new(CompositeInductiveDecisionVariable)

type CompositeInductiveDecisionVariable struct {
	variable.SimpleDecisionVariable

	weightedVariables map[variable.UndoableDecisionVariable]float64

	variableScales    map[variable.UndoableDecisionVariable]float64
	scaleVectorLength float64

	variable.ContainedDecisionVariableObservers
}

func (v *CompositeInductiveDecisionVariable) Initialise() *CompositeInductiveDecisionVariable {
	v.weightedVariables = make(map[variable.UndoableDecisionVariable]float64, 0)
	v.variableScales = make(map[variable.UndoableDecisionVariable]float64, 0)
	return v
}

func (v *CompositeInductiveDecisionVariable) WithName(name string) *CompositeInductiveDecisionVariable {
	v.SetName(name)
	return v
}

func (v *CompositeInductiveDecisionVariable) WithWeightedVariable(variable variable.UndoableDecisionVariable, weight float64) *CompositeInductiveDecisionVariable {
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

	v.calculateScalingVector()
	return v, nil
}

func (v *CompositeInductiveDecisionVariable) calculateScalingVector() {
	// https://en.wikipedia.org/wiki/Feature_scaling#Scaling_to_unit_length
	v.scaleVectorLength = v.vectorLengthOfVariableValues()
	for variable := range v.weightedVariables {
		variableValue := variable.Value()
		v.variableScales[variable] = variableValue / v.scaleVectorLength
	}
}

func (v *CompositeInductiveDecisionVariable) vectorLengthOfVariableValues() float64 {
	var summedSquares float64
	for variable := range v.weightedVariables {
		summedSquares += math.Pow(variable.Value(), 2)
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
	numberOfVariables := float64(len(v.variableScales))
	value := float64(0)
	for variable, scale := range v.variableScales {
		variableValue := variable.Value()
		scaledValue := variableValue / scale
		weight := v.weightedVariables[variable]
		value += scaledValue * weight * numberOfVariables
	}
	return value / v.scaleVectorLength
}

func (v *CompositeInductiveDecisionVariable) SetValue(value float64) {
	// Deliberately does nothing
}

func (v *CompositeInductiveDecisionVariable) UndoableValue() float64 {
	numberOfVariables := float64(len(v.variableScales))
	value := float64(0)
	for variable, scale := range v.variableScales {
		variableValue := variable.UndoableValue()
		scaledValue := variableValue / scale
		weight := v.weightedVariables[variable]
		value += scaledValue * weight * numberOfVariables
	}
	return value / v.scaleVectorLength
}

func (v *CompositeInductiveDecisionVariable) SetUndoableValue(value float64) {
	// Deliberately does nothing
}

func (v *CompositeInductiveDecisionVariable) ApplyDoneValue() {
	for variable := range v.weightedVariables {
		variable.ApplyDoneValue()
	}
	v.NotifyObservers()
}

func (v *CompositeInductiveDecisionVariable) ApplyUndoneValue() {
	for variable := range v.weightedVariables {
		variable.ApplyUndoneValue()
	}
	v.NotifyObservers()
}

func (v *CompositeInductiveDecisionVariable) DifferenceInValues() float64 {
	return v.UndoableValue() - v.Value()
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variableOld.
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
	return json.Marshal(variable.MakeEncodeable(v))
}
