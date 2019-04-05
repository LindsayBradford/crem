// Copyright (c) 2019 Australian Rivers Institute.

package solution

import (
	"sort"
	"strings"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
)

func NewSolution(id string) *Solution {
	newSolution := new(Solution)

	newSolution.Id = id
	newSolution.DecisionVariables = make(variable.EncodeableDecisionVariables, 0)
	newSolution.ActiveManagementActions = make(map[PlanningUnitId]ManagementActions, 0)

	return newSolution
}

type PlanningUnitId string
type PlanningUnitIds []PlanningUnitId

type ManagementActionType string

type ManagementActions []ManagementActionType

func (m ManagementActions) Len() int {
	return len(m)
}

func (m ManagementActions) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m ManagementActions) Less(i, j int) bool {
	return m[i] < m[j]
}

type Solution struct {
	Id                      string
	DecisionVariables       variable.EncodeableDecisionVariables
	PlanningUnits           PlanningUnitIds `json:"-"`
	ActiveManagementActions map[PlanningUnitId]ManagementActions
}

func (s Solution) ActiveActionsAsStrings() []string {
	actionList := make(ManagementActions, 0)

	entryAdded := make(map[ManagementActionType]bool, 0)
	for _, actions := range s.ActiveManagementActions {
		for _, action := range actions {
			if _, hasEntry := entryAdded[action]; !hasEntry {
				actionList = append(actionList, action)
				entryAdded[action] = true
			}
		}
	}

	sort.Sort(actionList)

	return actionsToStrings(actionList)
}

func actionsToStrings(actionList ManagementActions) []string {
	stringList := make([]string, len(actionList))
	for i, action := range actionList {
		stringList[i] = string(action)
	}
	return stringList
}

func (s Solution) FileNameSafeId() string {
	safeId := strings.Replace(s.Id, " ", "", -1)
	safeId = strings.Replace(safeId, "/", "_of_", -1)
	return safeId
}
