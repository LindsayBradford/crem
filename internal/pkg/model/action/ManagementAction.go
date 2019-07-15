// Copyright (c) 2019 Australian Rivers Institute.

// action package contains interfaces and behaviour for general-purpose "management actions" that change a model's
// decision variables based on their activation status.
package action

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

// ManagementActionType identifies a set of management actions that make the same kind of change to relevant
// model decision variables.
type ManagementActionType string

func (t ManagementActionType) String() string {
	return string(t)
}

// ModelVariableName identifies a particular model variable value that a management action modifies when its
// activation status changes.
type ModelVariableName string

// ManagementAction defines a general interface for the implementation of management actions.
type ManagementAction interface {
	// PlanningUnit returns the identifier of the planning unit in which management action is spatially located.
	PlanningUnit() planningunit.Id

	// Type identifies the ManagementActionType of a particular management action
	Type() ManagementActionType

	// IsActive reports whether a management action is active (true) or not (false).
	IsActive() bool

	// ModelVariableName reports tha value of the model variableName stored with the management action.
	ModelVariableValue(variableName ModelVariableName) float64

	// Subscribe allows a number of implementations of Observer to subscribe fur updates (changes in activation state).
	Subscribe(observers ...Observer)

	// InitialisingActivation activates an inactive management action, triggering Observer method
	// ObserveActionInitialising callbacks.
	InitialisingActivation()

	// ToggleActivation activates an inactive ManagementAction and vice-versa,
	// triggering Observer method ObserveAction callbacks.
	ToggleActivation()

	// ToggleActivationUnobserved activates an inactive ManagementAction and vice-versa, without triggering any
	// Observer method callbacks. Expected to be called when undoing a change that observers shouldn't react to.
	ToggleActivationUnobserved()

	// SetActivation activates an the ManagementAction as per the value supplied,
	SetActivation(value bool)

	// SetActivation activates an the ManagementAction as per the value supplied, without triggering any
	//  Observer method callbacks. Expected to be called when undoing a change that observers shouldn't react to.
	SetActivationUnobserved(value bool)
}
