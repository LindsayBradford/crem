package annealers

import (
	"sort"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
)

type SolutionBuilder struct {
	id    string
	model model.Model

	solution *solution.Solution
}

func (sb *SolutionBuilder) WithId(id string) *SolutionBuilder {
	sb.id = id
	return sb
}

func (sb *SolutionBuilder) ForModel(model model.Model) *SolutionBuilder {
	sb.model = model
	return sb
}

func (sb *SolutionBuilder) Build() *solution.Solution {
	sb.solution = solution.NewSolution(sb.id)

	sb.addDecisionVariables()
	sb.addPlanningUnits()
	sb.addPlanningUnitManagementActionMap()

	return sb.solution
}

func (sb *SolutionBuilder) addDecisionVariables() {
	if sb.model.DecisionVariables() == nil {
		return
	}

	solutionVariables := make(variable.EncodeableDecisionVariables, 0)

	for _, rawVariable := range *sb.model.DecisionVariables() {
		solutionVariables = append(solutionVariables, variable.MakeEncodeable(rawVariable))
	}

	sort.Sort(solutionVariables)
	sb.solution.DecisionVariables = solutionVariables
}

func (sb *SolutionBuilder) addPlanningUnits() {
	if sb.model.PlanningUnits() == nil {
		return
	}

	sb.solution.PlanningUnits = sb.model.PlanningUnits()
}

func (sb *SolutionBuilder) addPlanningUnitManagementActionMap() {
	for _, action := range sb.model.ActiveManagementActions() {
		planningUnit := action.PlanningUnit()
		actionType := solution.ManagementActionType(action.Type())
		sb.solution.ActiveManagementActions[planningUnit] = append(sb.solution.ActiveManagementActions[planningUnit], actionType)
	}
}
