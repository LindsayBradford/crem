// Copyright (c) 2019 Australian Rivers Institute.

package variableOld

import (
	"github.com/pkg/errors"
)

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

// SetValue finds the variableOld with supplied name in its collection, and sets its Value appropriately.
// If the collection has no variableOld for the supplied name, it panics.
func (vs *InductiveDecisionVariables) SetValue(name string, value float64) {
	if variable, isPresent := vs.asMap()[name]; isPresent {
		variable.SetValue(value)
		return
	}
	panic(variableMissing(name))
}

func variableMissing(name string) error {
	return errors.New("decision variableOld [" + name + "] does not exist.")
}

// Variable returns a pointer to the variableOld in its collection with the supplied name.
// If the collection has no variableOld for the supplied name, it panics.
func (vs *InductiveDecisionVariables) Variable(name string) InductiveDecisionVariable {
	if variable, isPresent := vs.asMap()[name]; isPresent {
		return variable
	}
	panic(variableMissing(name))
}

// Value returns the Value of the variableOld in its collection with the supplied name.
// If the collection has no variableOld for the supplied name, it panics.
func (vs *InductiveDecisionVariables) Value(name string) float64 {
	if variable, isPresent := vs.asMap()[name]; isPresent {
		return variable.Value()
	}
	panic(variableMissing(name))
}

// DifferenceInValues reports the difference in values of the variableOld in its collection with the supplied name.
func (vs *InductiveDecisionVariables) DifferenceInValues(variableName string) float64 {
	decisionVariable := vs.Variable(variableName)
	return decisionVariable.DifferenceInValues()
}

// AcceptAll accepts the inductive Value of all the BaseInductiveDecisionVariable instances in its collection.
func (vs *InductiveDecisionVariables) AcceptAll() {
	for _, variable := range vs.asMap() {
		variable.AcceptInductiveValue()
	}
}

// RejectAll rejects the inductive Value of all the BaseInductiveDecisionVariable instances in its collection.
func (vs *InductiveDecisionVariables) RejectAll() {
	for _, variable := range vs.asMap() {
		variable.RejectInductiveValue()
	}
}
