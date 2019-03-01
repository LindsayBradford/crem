// Copyright (c) 2019 Australian Rivers Institute.

package variable

type VolatileDecisionVariable struct {
	actual    SimpleDecisionVariable
	temporary SimpleDecisionVariable

	ContainedDecisionVariableObservers
}

func (v *VolatileDecisionVariable) Name() string {
	return v.actual.name
}

func (v *VolatileDecisionVariable) SetName(name string) {
	v.actual.name = name
}

func (v *VolatileDecisionVariable) Value() float64 {
	return v.actual.value
}

func (v *VolatileDecisionVariable) SetValue(value float64) {
	v.actual.value = value
	v.temporary.value = value
}

func (v *VolatileDecisionVariable) TemporaryValue() float64 {
	return v.temporary.value
}

func (v *VolatileDecisionVariable) ChangeInValue() float64 {
	return v.TemporaryValue() - v.Value()
}

func (v *VolatileDecisionVariable) SetTemporaryValue(value float64) {
	v.temporary.value = value
}

func (v *VolatileDecisionVariable) Accept() {
	v.actual.value = v.temporary.value
	v.NotifyObservers()
}

func (v *VolatileDecisionVariable) Revert() {
	v.temporary.value = v.actual.value
	v.NotifyObservers()
}

func (v *VolatileDecisionVariable) NotifyObservers() {
	for _, observer := range v.observers {
		observer.ObserveDecisionVariable(v)
	}
}
