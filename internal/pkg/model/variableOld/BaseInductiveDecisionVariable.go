// Copyright (c) 2019 Australian Rivers Institute.

package variableOld

import (
	"encoding/json"

	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
)

var _ variable.UndoableDecisionVariable = new(BaseInductiveDecisionVariable)

// BaseInductiveDecisionVariable offers a simple implementation of the InductiveDecisionVariable interface with
// the expectation that specific decisions variables will embed this struct to make use of typical
// InductiveDecisionVariable behaviour.
type BaseInductiveDecisionVariable struct {
	variable.SimpleDecisionVariable

	actualValue    float64
	inductiveValue float64

	variable.ContainedDecisionVariableObservers
}

func (v *BaseInductiveDecisionVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(variable.MakeEncodeable(v))
}

func (v *BaseInductiveDecisionVariable) Value() float64 {
	return v.actualValue
}

func (v *BaseInductiveDecisionVariable) SetValue(value float64) {
	v.actualValue = value
	v.inductiveValue = value
}

func (v *BaseInductiveDecisionVariable) UndoableValue() float64 {
	return v.inductiveValue
}

func (v *BaseInductiveDecisionVariable) DifferenceInValues() float64 {
	return v.UndoableValue() - v.Value()
}

func (v *BaseInductiveDecisionVariable) SetUndoableValue(value float64) {
	v.inductiveValue = value
}

func (v *BaseInductiveDecisionVariable) ApplyDoneValue() {
	if v.actualValue == v.inductiveValue {
		return
	}
	v.actualValue = v.inductiveValue
	v.NotifyObservers()
}

func (v *BaseInductiveDecisionVariable) ApplyUndoneValue() {
	if v.actualValue == v.inductiveValue {
		return
	}
	v.inductiveValue = v.actualValue
	v.NotifyObservers()
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variableOld.
func (v *BaseInductiveDecisionVariable) NotifyObservers() {
	for _, observer := range v.Observers() {
		observer.ObserveDecisionVariable(v)
	}
}

func (v *BaseInductiveDecisionVariable) ObserveAction(action action.ManagementAction) {
	// deliberately does nothing
}

func (v *BaseInductiveDecisionVariable) ObserveActionInitialising(action action.ManagementAction) {
	// deliberately does nothing
}
