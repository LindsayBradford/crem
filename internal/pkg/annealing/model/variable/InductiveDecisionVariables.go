// Copyright (c) 2019 Australian Rivers Institute.

package variable

import "github.com/pkg/errors"

func NewInductiveDecisionVariables() InductiveDecisionVariables {
	return make(InductiveDecisionVariables, 1)
}

type InductiveDecisionVariables map[string]*InductiveDecisionVariable

func (vs *InductiveDecisionVariables) Add(newVariables ...*InductiveDecisionVariable) {
	for _, newVariable := range newVariables {
		vs.asMap()[newVariable.Name()] = newVariable
	}
}

func (vs *InductiveDecisionVariables) NewForName(name string) {
	newVariable := new(InductiveDecisionVariable)
	newVariable.SetName(name)
	vs.asMap()[name] = newVariable
}

func (vs *InductiveDecisionVariables) asMap() InductiveDecisionVariables {
	return *vs
}

func (vs *InductiveDecisionVariables) SetValue(name string, value float64) {
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

func (vs *InductiveDecisionVariables) Variable(name string) *InductiveDecisionVariable {
	foundEntry, present := vs.asMap()[name]
	if !present {
		panic(variableMissing(name))
	}
	return foundEntry
}

func (vs *InductiveDecisionVariables) Value(name string) float64 {
	foundEntry, present := vs.asMap()[name]
	if !present {
		panic(variableMissing(name))
	}
	return foundEntry.Value()
}

func (vs *InductiveDecisionVariables) Accept() {
	for _, variable := range vs.asMap() {
		variable.AcceptInductiveValue()
	}
}

func (vs *InductiveDecisionVariables) Revert() {
	for _, variable := range vs.asMap() {
		variable.RejectInductiveValue()
	}
}
