// Copyright (c) 2019 Australian Rivers Institute.

package variable

import "github.com/pkg/errors"

func NewVolatileDecisionVariables() VolatileDecisionVariables {
	return make(VolatileDecisionVariables, 1)
}

type VolatileDecisionVariables map[string]*VolatileDecisionVariable

func (vs *VolatileDecisionVariables) Add(newVariables ...*VolatileDecisionVariable) {
	for _, newVariable := range newVariables {
		vs.asMap()[newVariable.Name()] = newVariable
	}
}

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

func (vs *VolatileDecisionVariables) Accept() {
	for _, variable := range vs.asMap() {
		variable.Accept()
	}
}

func (vs *VolatileDecisionVariables) Revert() {
	for _, variable := range vs.asMap() {
		variable.Revert()
	}
}
