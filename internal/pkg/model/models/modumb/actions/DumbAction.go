// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

var _ action.ManagementAction = new(DumbAction)

func New() *DumbAction {
	newAction := new(DumbAction).WithType(DumbActionType)
	return newAction
}

type DumbAction struct {
	action.SimpleManagementAction
}

const DumbActionType action.ManagementActionType = "DumbAction"

func (r *DumbAction) WithType(value action.ManagementActionType) *DumbAction {
	r.SimpleManagementAction.WithType(value)
	return r
}

func (r *DumbAction) WithPlanningUnit(planningUnit planningunit.Id) *DumbAction {
	r.SimpleManagementAction.WithPlanningUnit(planningUnit)
	return r
}

func (r *DumbAction) WithObjectiveValue(objectiveName action.ModelVariableName, objectiveValue float64) *DumbAction {
	r.SimpleManagementAction.WithVariable(objectiveName, objectiveValue)
	return r
}
