// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"github.com/LindsayBradford/crem/pkg/name"
)

var _ InductiveDecisionVariable = new(CompositeInductiveDecisionVariable)

type CompositeInductiveDecisionVariable struct {
	name.ContainedName

	composedVariables []InductiveDecisionVariable

	ContainedDecisionVariableObservers
}

func (v *CompositeInductiveDecisionVariable) Initialise() *CompositeInductiveDecisionVariable {
	v.composedVariables = make([]InductiveDecisionVariable, 0)
	return v
}

func (v *CompositeInductiveDecisionVariable) WithName(name string) *CompositeInductiveDecisionVariable {
	v.SetName(name)
	return v
}

func (v *CompositeInductiveDecisionVariable) ComposedOf(variables ...InductiveDecisionVariable) *CompositeInductiveDecisionVariable {
	v.composedVariables = append(v.composedVariables, variables...)
	return v
}

func (v *CompositeInductiveDecisionVariable) Value() float64 {
	value := float64(0)
	for _, variable := range v.composedVariables {
		value = value + variable.Value()
	}
	return value
}

func (v *CompositeInductiveDecisionVariable) SetValue(value float64) {
	// Deliberately does nothing
}

func (v *CompositeInductiveDecisionVariable) InductiveValue() float64 {
	inductiveValue := float64(0)
	for _, variable := range v.composedVariables {
		inductiveValue = inductiveValue + variable.InductiveValue()
	}
	return inductiveValue
}

func (v *CompositeInductiveDecisionVariable) SetInductiveValue(value float64) {
	// Deliberately does nothing
}

func (v *CompositeInductiveDecisionVariable) AcceptInductiveValue() {
	for _, variable := range v.composedVariables {
		variable.AcceptInductiveValue()
	}
	v.NotifyObservers()
}

func (v *CompositeInductiveDecisionVariable) RejectInductiveValue() {
	for _, variable := range v.composedVariables {
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
