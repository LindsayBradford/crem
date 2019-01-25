// Copyright (c) 2019 Australian Rivers Institute.

package action

type ManagementActionType string
type ModelVariableName string

type ManagementAction interface {
	PlanningUnit() string
	Type() ManagementActionType
	IsActive() bool
	ModelVariableValue(variableName ModelVariableName) float64
}

var _ ManagementAction = NewSimpleManagementAction()

const (
	active   = true
	inactive = false
)

func NewSimpleManagementAction() *SimpleManagementAction {
	newAction := new(SimpleManagementAction)
	newAction.variables = make(map[ModelVariableName]float64, 0)
	return newAction
}

type SimpleManagementAction struct {
	planningUnit string
	actionType   ManagementActionType
	isActive     bool

	variables map[ModelVariableName]float64
}

func (sma *SimpleManagementAction) SetPlanningUnit(planningUnit string) {
	sma.planningUnit = planningUnit
}

func (sma *SimpleManagementAction) PlanningUnit() string {
	return sma.planningUnit
}

func (sma *SimpleManagementAction) SetType(actionType ManagementActionType) {
	sma.actionType = actionType
}

func (sma *SimpleManagementAction) Type() ManagementActionType {
	return sma.actionType
}

func (sma *SimpleManagementAction) ToggleActivation() {
	sma.isActive = !sma.isActive
}

func (sma *SimpleManagementAction) Activate() {
	sma.isActive = active
}

func (sma *SimpleManagementAction) Deactivate() {
	sma.isActive = inactive
}

func (sma *SimpleManagementAction) IsActive() bool {
	return sma.isActive
}

func (sma *SimpleManagementAction) SetModelVariable(variableName ModelVariableName, value float64) {
	sma.variables[variableName] = value
}

func (sma *SimpleManagementAction) ModelVariableValue(variableName ModelVariableName) float64 {
	return sma.variables[variableName]
}
