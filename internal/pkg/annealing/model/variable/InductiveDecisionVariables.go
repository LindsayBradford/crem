// Copyright (c) 2019 Australian Rivers Institute.

package variable

import "github.com/pkg/errors"

func NewInductiveDecisionVariables() InductiveDecisionVariables {
	return make(InductiveDecisionVariables, 1)
}

// InductiveDecisionVariables offers up a name-indexed collection of InductiveDecisionVariable instances, along with
// convenience methods for the collection's management.  It is typically expected that a model would contain only a
// single instance of InductiveDecisionVariables to house all of its decision variables.
type InductiveDecisionVariables map[string]InductiveDecisionVariable

// Adds a number of InductiveDecisionVariables to the collection
func (vs *InductiveDecisionVariables) Add(newVariables ...InductiveDecisionVariable) {
	for _, newVariable := range newVariables {
		vs.asMap()[newVariable.Name()] = newVariable
	}
}

// NewForName creates and adds to its collection, a new BaseInductiveDecisionVariable with the supplied name.
func (vs *InductiveDecisionVariables) NewForName(name string) {
	newVariable := new(BaseInductiveDecisionVariable)
	newVariable.SetName(name)
	vs.asMap()[name] = newVariable
}

func (vs *InductiveDecisionVariables) asMap() InductiveDecisionVariables {
	return *vs
}

// SetValue finds the variable with supplied name in its collection, and sets its value appropriately.
// If the collection has no variable for the supplied name, it panics.
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

// Variable returns a pointer to the variable in its collection with the supplied name.
// If the collection has no variable for the supplied name, it panics.
func (vs *InductiveDecisionVariables) Variable(name string) InductiveDecisionVariable {
	foundEntry, present := vs.asMap()[name]
	if !present {
		panic(variableMissing(name))
	}
	return foundEntry
}

// Value returns the value of the variable in its collection with the supplied name.
// If the collection has no variable for the supplied name, it panics.
func (vs *InductiveDecisionVariables) Value(name string) float64 {
	foundEntry, present := vs.asMap()[name]
	if !present {
		panic(variableMissing(name))
	}
	return foundEntry.Value()
}

// AcceptAll accepts the inductive value of all the BaseInductiveDecisionVariable instances in its collection.
func (vs *InductiveDecisionVariables) AcceptAll() {
	for _, variable := range vs.asMap() {
		variable.AcceptInductiveValue()
	}
}

// RejectAll rejects the inductive value of all the BaseInductiveDecisionVariable instances in its collection.
func (vs *InductiveDecisionVariables) RejectAll() {
	for _, variable := range vs.asMap() {
		variable.RejectInductiveValue()
	}
}
