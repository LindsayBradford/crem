// Copyright (c) 2019 Australian Rivers Institute.

package variable

import "github.com/LindsayBradford/crem/pkg/name"

// InductiveDecisionVariable is a DecisionVariable that allows an 'inductive' value to be temporarily stored
// and retrieved for the decision variable (typically based based on some management action).
// The induced value does not become the actual value for the decision variable without being explicitly accepted.
// The induced value can also be rejected, which sees it revert to the actual value of the variable.
type InductiveDecisionVariable struct {
	name.ContainedName

	actualValue    float64
	inductiveValue float64

	ContainedDecisionVariableObservers
}

func (v *InductiveDecisionVariable) Value() float64 {
	return v.actualValue
}

func (v *InductiveDecisionVariable) SetValue(value float64) {
	v.actualValue = value
	v.inductiveValue = value
}

func (v *InductiveDecisionVariable) InductiveValue() float64 {
	return v.inductiveValue
}

func (v *InductiveDecisionVariable) DifferenceInValues() float64 {
	return v.InductiveValue() - v.Value()
}

func (v *InductiveDecisionVariable) SetInductiveValue(value float64) {
	v.inductiveValue = value
}

func (v *InductiveDecisionVariable) AcceptInductiveValue() {
	v.actualValue = v.inductiveValue
	v.NotifyObservers()
}

func (v *InductiveDecisionVariable) RejectInductiveValue() {
	v.inductiveValue = v.actualValue
	v.NotifyObservers()
}

func (v *InductiveDecisionVariable) NotifyObservers() {
	for _, observer := range v.observers {
		observer.ObserveDecisionVariable(v)
	}
}
