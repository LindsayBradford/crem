// Copyright (c) 2019 Australian Rivers Institute.

package solution

import (
	"strings"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
)

func NewSolution(id string) *Solution {
	newSolution := new(Solution)

	newSolution.Id = id
	newSolution.DecisionVariables = make(variable.EncodeableDecisionVariables, 0)
	newSolution.PlanningUnitManagementActionsMap = make(map[PlanningUnitId]ManagementActions, 0)

	return newSolution
}

type PlanningUnitId string
type ManagementActionType string

type ManagementActions []ManagementActionType

type Solution struct {
	Id                               string
	DecisionVariables                variable.EncodeableDecisionVariables
	PlanningUnitManagementActionsMap map[PlanningUnitId]ManagementActions
}

func (s Solution) FileNameSafeId() string {
	safeId := strings.Replace(s.Id, " ", "", -1)
	safeId = strings.Replace(safeId, "/", "_of_", -1)
	return safeId
}
