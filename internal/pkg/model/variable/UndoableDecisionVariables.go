// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"github.com/pkg/errors"
)

func NewUndoableDecisionVariables() UndoableDecisionVariables {
	return make(UndoableDecisionVariables, 0)
}

// UndoableDecisionVariables offers up a name-indexed collection of UndoableDecisionVariable instances, along with
// convenience methods for the collection's management.  It is typically expected that a model would contain only a
// single instance of UndoableDecisionVariables to house all of its decision variables.
type UndoableDecisionVariables []UndoableDecisionVariable

// Adds a number of UndoableDecisionVariables to the collection
func (vs *UndoableDecisionVariables) Add(newVariables ...UndoableDecisionVariable) {
	*vs = append(*vs, newVariables...)
}

// NewForName creates and adds to its colle ction, a new BaseInductiveDecisionVariable with the supplied name.
func (vs *UndoableDecisionVariables) NewForName(name string) {
	newVariable := new(SimpleUndoableDecisionVariable)
	newVariable.SetName(name)

	*vs = append(*vs, newVariable)
}

// SetValue finds the variableOld with supplied name in its collection, and sets its Value appropriately.
// If the collection has no variableOld for the supplied name, it panics.
func (vs *UndoableDecisionVariables) SetValue(name string, value float64) {
	vs.find(name).SetValue(value)
}

func variableMissing(name string) error {
	return errors.New("decision variable [" + name + "] does not exist.")
}

// Variable returns a pointer to the variableOld in its collection with the supplied name.
// If the collection has no variableOld for the supplied name, it panics.
func (vs *UndoableDecisionVariables) Variable(name string) UndoableDecisionVariable {
	return vs.find(name)
}

func (vs *UndoableDecisionVariables) find(name string) UndoableDecisionVariable {
	for _, variable := range *vs {
		if variable.Name() == name {
			return variable
		}
	}
	panic(variableMissing(name))
}

// Value returns the Value of the variableOld in its collection with the supplied name.
// If the collection has no variableOld for the supplied name, it panics.
func (vs *UndoableDecisionVariables) Value(name string) float64 {
	return vs.find(name).Value()
}

// DifferenceInValues reports the difference in values of the variableOld in its collection with the supplied name.
func (vs *UndoableDecisionVariables) DifferenceInValues(variableName string) float64 {
	decisionVariable := vs.Variable(variableName)
	return decisionVariable.DifferenceInValues()
}

// AcceptAll accepts the inductive Value of all the BaseInductiveDecisionVariable instances in its collection.
func (vs *UndoableDecisionVariables) AcceptAll() {
	for _, variable := range *vs {
		variable.ApplyDoneValue()
	}
}

// RejectAll rejects the inductive Value of all the BaseInductiveDecisionVariable instances in its collection.
func (vs *UndoableDecisionVariables) RejectAll() {
	for _, variable := range *vs {
		variable.ApplyUndoneValue()
	}
}
