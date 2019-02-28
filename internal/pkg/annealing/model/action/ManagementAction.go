// Copyright (c) 2019 Australian Rivers Institute.

package action

type ManagementActionType string
type ModelVariableName string

type ManagementAction interface {
	PlanningUnit() string
	Type() ManagementActionType
	IsActive() bool
	ModelVariableValue(variableName ModelVariableName) float64

	Subscribe(observers ...Observer)

	ToggleInitialisingActivation()
	ToggleActivation()
	ToggleActivationUnobserved()
}
