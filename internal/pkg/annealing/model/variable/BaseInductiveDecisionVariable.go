// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"encoding/json"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/action"
	"github.com/LindsayBradford/crem/pkg/name"
)

var _ InductiveDecisionVariable = new(BaseInductiveDecisionVariable)

// BaseInductiveDecisionVariable offers a simple implementation of the InductiveDecisionVariable interface with
// the expectation that specific decisions variables will embed this struct to make use of typical
// InductiveDecisionVariable behaviour.
type BaseInductiveDecisionVariable struct {
	name.NameContainer

	actualValue    float64
	inductiveValue float64

	unitOfMeasure UnitOfMeasure

	ContainedDecisionVariableObservers
	ContainedUnitOfMeasure
}

func (v *BaseInductiveDecisionVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(MakeEncodeable(v))
}

func (v *BaseInductiveDecisionVariable) Value() float64 {
	return v.actualValue
}

func (v *BaseInductiveDecisionVariable) SetValue(value float64) {
	v.actualValue = value
	v.inductiveValue = value
}

func (v *BaseInductiveDecisionVariable) InductiveValue() float64 {
	return v.inductiveValue
}

func (v *BaseInductiveDecisionVariable) DifferenceInValues() float64 {
	return v.InductiveValue() - v.Value()
}

func (v *BaseInductiveDecisionVariable) SetInductiveValue(value float64) {
	v.inductiveValue = value
}

func (v *BaseInductiveDecisionVariable) AcceptInductiveValue() {
	if v.actualValue == v.inductiveValue {
		return
	}
	v.actualValue = v.inductiveValue
	v.NotifyObservers()
}

func (v *BaseInductiveDecisionVariable) RejectInductiveValue() {
	if v.actualValue == v.inductiveValue {
		return
	}
	v.inductiveValue = v.actualValue
	v.NotifyObservers()
}

// NotifyObservers allows structs embedding a BaseInductiveDecisionVariable to trigger a notification of change
// to any observers watching for state changes to the variable.
func (v *BaseInductiveDecisionVariable) NotifyObservers() {
	for _, observer := range v.observers {
		observer.ObserveDecisionVariable(v)
	}
}

func (v *BaseInductiveDecisionVariable) ObserveAction(action action.ManagementAction) {
	// deliberately does nothing
}

func (v *BaseInductiveDecisionVariable) ObserveActionInitialising(action action.ManagementAction) {
	// deliberately does nothing
}
