// Copyright (c) 2019 Australian Rivers Institute.

package action

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"math"
)

// NullManagementAction is a stateless/neutral implementation of the ManagementAction interface, useful for testing
// scenarios where action state is irrelevant or within scenarios where a Null Pattern for actions makes sense
// (https://en.wikipedia.org/wiki/Null_object_pattern).
var NullManagementAction ManagementAction = new(Null)

type Null struct{}

const NullPlanningUnitId planningunit.Id = math.MaxUint64

func (a *Null) PlanningUnit() planningunit.Id { return NullPlanningUnitId }

const NullManagementActionType ManagementActionType = "NullType"

func (a *Null) Type() ManagementActionType                                { return NullManagementActionType }
func (a *Null) IsActive() bool                                            { return false }
func (a *Null) ModelVariableValue(variableName ModelVariableName) float64 { return 0 }
func (a *Null) Subscribe(observers ...Observer)                           {}
func (a *Null) InitialisingActivation()                                   {}
func (a *Null) ToggleActivation()                                         {}
func (a *Null) ToggleActivationUnobserved()                               {}
func (a *Null) SetActivation(value bool)                                  {}
func (a *Null) SetActivationUnobserved(value bool)                        {}
