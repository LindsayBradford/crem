// Copyright (c) 2019 Australian Rivers Institute.

package action

var _ ManagementAction = new(SimpleManagementAction)

const (
	active   = true
	inactive = false
)

type SimpleManagementAction struct {
	planningUnit string
	actionType   ManagementActionType
	isActive     bool

	variables map[ModelVariableName]float64
	observers []Observer
}

func (sma *SimpleManagementAction) WithPlanningUnit(planningUnit string) *SimpleManagementAction {
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

func (sma *SimpleManagementAction) PlanningUnit() string {
	return sma.planningUnit
}

func (sma *SimpleManagementAction) Type() ManagementActionType {
	return sma.actionType
}

func (sma *SimpleManagementAction) ToggleActivation() {
	sma.ToggleActivationUnobserved()
	sma.notifyObservers()
}

func (sma *SimpleManagementAction) ToggleActivationUnobserved() {
	sma.isActive = !sma.isActive
}

func (sma *SimpleManagementAction) Activate() {
	if sma.isActive {
		return
	}
	sma.isActive = active
	sma.notifyObservers()
}

func (sma *SimpleManagementAction) Deactivate() {
	if !sma.isActive {
		return
	}
	sma.isActive = inactive
	sma.notifyObservers()
}

func (sma *SimpleManagementAction) IsActive() bool {
	return sma.isActive
}

func (sma *SimpleManagementAction) ModelVariableValue(variableName ModelVariableName) float64 {
	return sma.variables[variableName]
}

func (sma *SimpleManagementAction) Subscribe(observers ...Observer) {
	if sma.observers == nil {
		sma.observers = make([]Observer, 0)
	}

	for _, newObserver := range observers {
		sma.observers = append(sma.observers, newObserver)
	}
}

func (sma *SimpleManagementAction) notifyObservers() {
	for _, observer := range sma.observers {
		observer.Observe(sma)
	}
}
