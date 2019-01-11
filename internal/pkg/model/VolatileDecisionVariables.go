// Copyright (c) 2019 Australian Rivers Institute.

package model

type VolatileDecisionVariables map[string]*VolatileDecisionVariable

func NewVolatileDecisionVariables() VolatileDecisionVariables {
	return make(VolatileDecisionVariables, 1)
}

type VolatileDecisionVariable struct {
	actual    DecisionVariableImpl
	temporary DecisionVariableImpl
}

func (dvi *VolatileDecisionVariable) Name() string {
	return dvi.actual.name
}

func (dvi *VolatileDecisionVariable) SetName(name string) {
	dvi.actual.name = name
}

func (dvi *VolatileDecisionVariable) Value() float64 {
	return dvi.actual.value
}

func (dvi *VolatileDecisionVariable) SetValue(value float64) {
	dvi.actual.value = value
	dvi.temporary.value = value
}

func (dvi *VolatileDecisionVariable) TemporaryValue() float64 {
	return dvi.temporary.value
}

func (dvi *VolatileDecisionVariable) SetTemporaryValue(value float64) {
	dvi.temporary.value = value
}

func (dvi *VolatileDecisionVariable) Accept() {
	dvi.actual.value = dvi.temporary.value
}

func (dvi *VolatileDecisionVariable) Revert() {
	dvi.temporary.value = dvi.actual.value
}
