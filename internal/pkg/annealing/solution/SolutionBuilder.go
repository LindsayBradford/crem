// Copyright (c) 2020 Australian Rivers Institute.

package solution

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/archive"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"sort"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
)

type SolutionBuilder struct {
	id              string
	model           model.Model
	compressedModel *archive.CompressedModelState

	solution *Solution
}

func (sb *SolutionBuilder) WithId(id string) *SolutionBuilder {
	sb.id = id
	return sb
}

func (sb *SolutionBuilder) ForModel(model model.Model) *SolutionBuilder {
	sb.model = model
	sb.compressedModel = new(archive.ModelCompressor).Compress(model)
	return sb
}

func (sb *SolutionBuilder) Build() *Solution {
	sb.solution = NewSolution(sb.id)
	sb.transferAttributes()
	sb.addDecisionVariables()
	sb.addPlanningUnits()
	sb.addPlanningUnitManagementActionMaps()

	return sb.solution
}

func (sb *SolutionBuilder) transferAttributes() {
	if attributeContainingModel, hasAttributes := sb.model.(attributes.Interface); hasAttributes {
		sb.solution.JoiningAttributes(attributeContainingModel.AllAttributes())
	}
}

func (sb *SolutionBuilder) addDecisionVariables() {
	if sb.model.NameMappedVariables() == nil {
		return
	}

	solutionVariables := make(variable.EncodeableDecisionVariables, 0)

	for _, rawVariable := range *sb.model.NameMappedVariables() {
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

func (sb *SolutionBuilder) addPlanningUnitManagementActionMaps() {
	sb.solution.EncodedActions = sb.compressedModel.Encoding()

	for _, action := range sb.model.ManagementActions() {
		planningUnit := action.PlanningUnit()
		actionType := ManagementActionType(action.Type())
		sb.solution.ManagementActions[actionType] = true
		switch action.IsActive() {
		case true:
			sb.solution.ActiveManagementActions[planningUnit] =
				append(sb.solution.ActiveManagementActions[planningUnit], actionType)
		case false:
			sb.solution.InactiveManagementActions[planningUnit] =
				append(sb.solution.InactiveManagementActions[planningUnit], actionType)
		}
	}
}
