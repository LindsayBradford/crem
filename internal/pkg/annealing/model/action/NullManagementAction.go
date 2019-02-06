// Copyright (c) 2019 Australian Rivers Institute.

package action

var NullManagementAction ManagementAction = new(Null)

type Null struct{}

func (a *Null) PlanningUnit() string { return "null" }

const NullManagementActionType ManagementActionType = "NullType"

func (a *Null) Type() ManagementActionType                                { return NullManagementActionType }
func (a *Null) IsActive() bool                                            { return false }
func (a *Null) ModelVariableValue(variableName ModelVariableName) float64 { return 0 }
func (a *Null) Subscribe(observers ...Observer)                           {}
func (a *Null) ToggleActivation()                                         {}
func (a *Null) ToggleActivationUnobserved()                               {}
