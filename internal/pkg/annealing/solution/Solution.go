// Copyright (c) 2019 Australian Rivers Institute.

package solution

import (
	"github.com/LindsayBradford/crem/pkg/attributes"
)

func NewSolution(id string) *Solution {
	newSolution := new(Solution)

	newSolution.Id = id
	newSolution.DecisionVariables = make(attributes.Attributes, 0)
	newSolution.PlanningUnitManagementActionsMap = make(map[PlanningUnitId]ManagementActions, 0)

	return newSolution
}

type PlanningUnitId string
type ManagementActionType string

type ManagementActions []ManagementActionType

type Solution struct {
	Id                               string
	DecisionVariables                attributes.Attributes
	PlanningUnitManagementActionsMap map[PlanningUnitId]ManagementActions
}
