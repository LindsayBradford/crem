// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"github.com/pkg/errors"
)

func NewUndoableDecisionVariables() UndoableDecisionVariables {
	return make(UndoableDecisionVariables, 1)
}

// UndoableDecisionVariables offers up a name-indexed collection of UndoableDecisionVariable instances, along with
// convenience methods for the collection's management.  It is typically expected that a model would contain only a
// single instance of UndoableDecisionVariables to house all of its decision variables.
type UndoableDecisionVariables map[string]UndoableDecisionVariable

// Adds a number of UndoableDecisionVariables to the collection
func (vs *UndoableDecisionVariables) Add(newVariables ...UndoableDecisionVariable) {
	for _, newVariable := range newVariables {
		vs.asMap()[newVariable.Name()] = newVariable
	}
}

// NewForName creates and adds to its collection, a new BaseInductiveDecisionVariable with the supplied name.
func (vs *UndoableDecisionVariables) NewForName(name string) {
	newVariable := new(SimpleUndoableDecisionVariable)
	newVariable.SetName(name)
	vs.asMap()[name] = newVariable
}

func (vs *UndoableDecisionVariables) asMap() UndoableDecisionVariables {
	return *vs
}

// SetValue finds the variableOld with supplied name in its collection, and sets its Value appropriately.
// If the collection has no variableOld for the supplied name, it panics.
func (vs *UndoableDecisionVariables) SetValue(name string, value float64) {
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
func (vs *UndoableDecisionVariables) Variable(name string) UndoableDecisionVariable {
	if variable, isPresent := vs.asMap()[name]; isPresent {
		return variable
	}
	panic(variableMissing(name))
}

// Value returns the Value of the variableOld in its collection with the supplied name.
// If the collection has no variableOld for the supplied name, it panics.
func (vs *UndoableDecisionVariables) Value(name string) float64 {
	if variable, isPresent := vs.asMap()[name]; isPresent {
		return variable.Value()
	}
	panic(variableMissing(name))
}

// DifferenceInValues reports the difference in values of the variableOld in its collection with the supplied name.
func (vs *UndoableDecisionVariables) DifferenceInValues(variableName string) float64 {
	decisionVariable := vs.Variable(variableName)
	return decisionVariable.DifferenceInValues()
}

// AcceptAll accepts the inductive Value of all the BaseInductiveDecisionVariable instances in its collection.
func (vs *UndoableDecisionVariables) AcceptAll() {
	for _, variable := range vs.asMap() {
		variable.ApplyDoneValue()
	}
}

// RejectAll rejects the inductive Value of all the BaseInductiveDecisionVariable instances in its collection.
func (vs *UndoableDecisionVariables) RejectAll() {
	for _, variable := range vs.asMap() {
		variable.ApplyUndoneValue()
	}
}
