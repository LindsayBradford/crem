// Copyright (c) 2019 Australian Rivers Institute.

package model

import "github.com/pkg/errors"

func NewVolatileDecisionVariables() VolatileDecisionVariables {
	return make(VolatileDecisionVariables, 1)
}

type VolatileDecisionVariables map[string]*VolatileDecisionVariable

func (vs *VolatileDecisionVariables) NewForName(name string) {
	newVariable := new(VolatileDecisionVariable)
	newVariable.SetName(name)
	vs.asMap()[name] = newVariable
}

func (vs *VolatileDecisionVariables) asMap() VolatileDecisionVariables {
	return *vs
}

func (vs *VolatileDecisionVariables) SetValue(name string, value float64) {
	foundEntry, present := vs.asMap()[name]
	if !present {
		panic(variableMissing(name))
	}
	value = foundEntry.Value()
	foundEntry.SetValue(value)
}

func variableMissing(name string) error {
	return errors.New("decision variable[" + name + "] does not exist.")
}

func (vs *VolatileDecisionVariables) Variable(name string) *VolatileDecisionVariable {
	foundEntry, present := vs.asMap()[name]
	if !present {
		panic(variableMissing(name))
	}
	return foundEntry
}

func (vs *VolatileDecisionVariables) Value(name string) float64 {
	foundEntry, present := vs.asMap()[name]
	if !present {
		panic(variableMissing(name))
	}
	return foundEntry.Value()
}

type VolatileDecisionVariable struct {
	actual    DecisionVariableImpl
	temporary DecisionVariableImpl
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

func (v *VolatileDecisionVariable) SetTemporaryValue(value float64) {
	v.temporary.value = value
}

func (v *VolatileDecisionVariable) Accept() {
	v.actual.value = v.temporary.value
}

func (v *VolatileDecisionVariable) Revert() {
	v.temporary.value = v.actual.value
}
