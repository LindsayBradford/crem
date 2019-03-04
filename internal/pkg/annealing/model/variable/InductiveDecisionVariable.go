// Copyright (c) 2019 Australian Rivers Institute.

package variable

import "github.com/LindsayBradford/crem/pkg/name"

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
