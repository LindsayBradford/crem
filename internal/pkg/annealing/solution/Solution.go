// Copyright (c) 2019 Australian Rivers Institute.

package solution

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"

	"sort"
	"strings"
)

func NewSolution(id string) *Solution {
	newSolution := new(Solution)

	newSolution.Id = id
	newSolution.DecisionVariables = make(variableNew.EncodeableDecisionVariables, 0)

	newSolution.ManagementActions = make(map[ManagementActionType]bool, 0)
	newSolution.ActiveManagementActions = make(map[planningunit.Id]ManagementActions, 0)
	newSolution.InactiveManagementActions = make(map[planningunit.Id]ManagementActions, 0)

	return newSolution
}

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
	Id                        string
	DecisionVariables         variableNew.EncodeableDecisionVariables
	PlanningUnits             planningunit.Ids              `json:"-"`
	ManagementActions         map[ManagementActionType]bool `json:"-"`
	ActiveManagementActions   map[planningunit.Id]ManagementActions
	InactiveManagementActions map[planningunit.Id]ManagementActions `json:"-"`
}

func (s Solution) ActionsAsStrings() []string {
	actionList := make(ManagementActions, 0)

	for actionKey := range s.ManagementActions {
		actionList = append(actionList, actionKey)
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
