// Copyright (c) 2019 Australian Rivers Institute.

package variable

import "github.com/LindsayBradford/crem/pkg/name"

var _ InductiveDecisionVariable = new(BaseInductiveDecisionVariable)

type BaseInductiveDecisionVariable struct {
	name.ContainedName

	actualValue    float64
	inductiveValue float64

	ContainedDecisionVariableObservers
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
	v.actualValue = v.inductiveValue
	v.NotifyObservers()
}

func (v *BaseInductiveDecisionVariable) RejectInductiveValue() {
	v.inductiveValue = v.actualValue
	v.NotifyObservers()
}

func (v *BaseInductiveDecisionVariable) NotifyObservers() {
	for _, observer := range v.observers {
		observer.ObserveDecisionVariable(v)
	}
}
