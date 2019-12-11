// Copyright (c) 2019 Australian Rivers Institute.

package solution

import (
	"fmt"
	"sort"
	"strings"

	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/math"
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

func (s *Solution) MatchErrors(other *Solution) *compositeErrors.CompositeError {
	matchErrors := compositeErrors.New("Solution Match Errors")

	s.checkIds(other, matchErrors)
	s.checkDecisionVariables(other, matchErrors)

	if matchErrors.Size() > 0 {
		return matchErrors
	}
	return nil
}

func (s *Solution) checkIds(other *Solution, errors *compositeErrors.CompositeError) {
	if s.Id != other.Id {
		idError := fmt.Sprintf("Solutions have mismatching Ids [%s, %s]", s.Id, other.Id)
		errors.AddMessage(idError)
	}
}

func (s *Solution) checkDecisionVariables(other *Solution, errors *compositeErrors.CompositeError) {
	s.checkForMissingDecisionVariables(other, errors)
	s.checkForMismatchedDecisionVariableValues(other, errors)
	s.checkDecisionVariablesAreSumOfPlanningUnits(other, errors)
}

func (s *Solution) checkForMissingDecisionVariables(other *Solution, errors *compositeErrors.CompositeError) {
	variableSolutionMap := make(map[string]*Solution, 0)

	for _, variable := range s.DecisionVariables {
		variableSolutionMap[variable.Name] = s
	}

	for _, variable := range other.DecisionVariables {
		if variableSolutionMap[variable.Name] != nil {
			delete(variableSolutionMap, variable.Name)
		} else {
			variableSolutionMap[variable.Name] = other
		}
	}

	for variableName, solution := range variableSolutionMap {
		variableError := fmt.Sprintf("Only solution [%s] has variable [%s]", solution.Id, variableName)
		errors.AddMessage(variableError)
	}
}

func (s *Solution) checkForMismatchedDecisionVariableValues(other *Solution, errors *compositeErrors.CompositeError) {
	for _, myVariable := range s.DecisionVariables {
		for _, otherVariable := range other.DecisionVariables {
			if myVariable.Name == otherVariable.Name {
				if myVariable.Value != otherVariable.Value {
					variableError := fmt.Sprintf("variable [%s] has mismatching values [%f, %f]", myVariable.Name, myVariable.Value, otherVariable.Value)
					errors.AddMessage(variableError)
				}
			}
		}
	}
}

func (s *Solution) checkDecisionVariablesAreSumOfPlanningUnits(other *Solution, errors *compositeErrors.CompositeError) {
	for _, myVariable := range s.DecisionVariables {
		var planningUnitValues float64
		for _, planningUnitValue := range myVariable.ValuePerPlanningUnit {
			planningUnitValues += planningUnitValue.Value
		}
		precisionOfVariable := math.DerivePrecision(myVariable.Value)
		if myVariable.Value != math.RoundFloat(planningUnitValues, precisionOfVariable) {
			variableError := fmt.Sprintf("variable [%s] has value [%f], but sum of planning units is [%f]", myVariable.Name, myVariable.Value, planningUnitValues)
			errors.AddMessage(variableError)
		}
	}
}
