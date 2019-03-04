// Copyright (c) 2019 Australian Rivers Institute.

package variable

type InductiveDecisionVariable struct {
	actual    SimpleDecisionVariable
	inductive SimpleDecisionVariable

	ContainedDecisionVariableObservers
}

func (v *InductiveDecisionVariable) Name() string {
	return v.actual.name
}

func (v *InductiveDecisionVariable) SetName(name string) {
	v.actual.name = name
}

func (v *InductiveDecisionVariable) Value() float64 {
	return v.actual.value
}

func (v *InductiveDecisionVariable) SetValue(value float64) {
	v.actual.value = value
	v.inductive.value = value
}

func (v *InductiveDecisionVariable) InductiveValue() float64 {
	return v.inductive.value
}

func (v *InductiveDecisionVariable) DifferenceInValues() float64 {
	return v.InductiveValue() - v.Value()
}

func (v *InductiveDecisionVariable) SetInductiveValue(value float64) {
	v.inductive.value = value
}

func (v *InductiveDecisionVariable) AcceptInductiveValue() {
	v.actual.value = v.inductive.value
	v.NotifyObservers()
}

func (v *InductiveDecisionVariable) RejectInductiveValue() {
	v.inductive.value = v.actual.value
	v.NotifyObservers()
}

func (v *InductiveDecisionVariable) NotifyObservers() {
	for _, observer := range v.observers {
		observer.ObserveDecisionVariable(v)
	}
}
