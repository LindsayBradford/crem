// Copyright (c) 2019 Australian Rivers Institute.

package action

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

var _ ManagementAction = new(SimpleManagementAction)

// SimpleManagementAction is a basic, generally useful implementation of the ManagementAction interface, using a
// fluent interface for its action construction.
type SimpleManagementAction struct {
	planningUnit planningunit.Id
	actionType   ManagementActionType
	isActive     bool

	variables map[ModelVariableName]float64
	observers []Observer
}

func (sma *SimpleManagementAction) WithPlanningUnit(planningUnit planningunit.Id) *SimpleManagementAction {
	sma.planningUnit = planningUnit
	return sma
}

func (sma *SimpleManagementAction) WithType(actionType ManagementActionType) *SimpleManagementAction {
	sma.actionType = actionType
	return sma
}

func (sma *SimpleManagementAction) WithVariable(variableName ModelVariableName, value float64) *SimpleManagementAction {
	if sma.variables == nil {
		sma.variables = make(map[ModelVariableName]float64, 0)
	}
	sma.variables[variableName] = value

	return sma
}

func (sma *SimpleManagementAction) PlanningUnit() planningunit.Id {
	return sma.planningUnit
}

func (sma *SimpleManagementAction) Type() ManagementActionType {
	return sma.actionType
}

func (sma *SimpleManagementAction) InitialisingActivation() {
	if sma.isActive {
		return
	}
	sma.ToggleActivationUnobserved()
	sma.notifyInitialisingObservers()
}

func (sma *SimpleManagementAction) InitialisingDeactivation() {
	if !sma.isActive {
		return
	}
	sma.ToggleActivationUnobserved()
	sma.notifyInitialisingObservers()
}

func (sma *SimpleManagementAction) ToggleActivation() {
	sma.ToggleActivationUnobserved()
	sma.notifyObservers()
}

func (sma *SimpleManagementAction) ToggleActivationUnobserved() {
	sma.isActive = !sma.isActive
}

func (sma *SimpleManagementAction) SetActivation(value bool) {
	sma.isActive = value
	sma.notifyObservers()
}

func (sma *SimpleManagementAction) SetActivationUnobserved(value bool) {
	sma.isActive = value
}

func (sma *SimpleManagementAction) IsActive() bool {
	return sma.isActive
}

func (sma *SimpleManagementAction) ModelVariableValue(variableName ModelVariableName) float64 {
	return sma.variables[variableName]
}

func (sma *SimpleManagementAction) ModelVariableKeys() (keys []ModelVariableName) {
	for key := range sma.variables {
		keys = append(keys, key)
	}
	return
}

func (sma *SimpleManagementAction) Subscribe(observers ...Observer) {
	sma.observers = observers
	//if sma.observers == nil {
	//	sma.observers = make([]Observer, 0)
	//}
	//
	//for _, newObserver := range observers {
	//	sma.observers = append(sma.observers, newObserver)
	//}
}

func (sma *SimpleManagementAction) notifyObservers() {
	for _, observer := range sma.observers {
		observer.ObserveAction(sma)
	}
}

func (sma *SimpleManagementAction) notifyInitialisingObservers() {
	for _, observer := range sma.observers {
		observer.ObserveActionInitialising(sma)
	}
}
