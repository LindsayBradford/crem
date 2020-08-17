// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

type GullyRestorationGroup struct {
	sedimentContribution *GullySedimentContribution
	parameters           parameters.Parameters

	actionMap map[planningunit.Id]*GullyRestoration
	Container
}

func (g *GullyRestorationGroup) WithParameters(parameters parameters.Parameters) *GullyRestorationGroup {
	g.parameters = parameters
	return g
}

func (g *GullyRestorationGroup) WithGullyTable(gullyTable tables.CsvTable) *GullyRestorationGroup {
	g.sedimentContribution = new(GullySedimentContribution)
	g.sedimentContribution.Initialise(gullyTable, g.parameters)
	return g
}

func (g *GullyRestorationGroup) WithActionsTable(actionsTable tables.CsvTable) *GullyRestorationGroup {
	g.Container.WithSourceFilter(GullySource).WithActionsTable(actionsTable)
	return g
}

func (g *GullyRestorationGroup) ManagementActions() []action.ManagementAction {
	g.createManagementActions()
	actions := make([]action.ManagementAction, 0)
	for _, value := range g.actionMap {
		actions = append(actions, value)
	}
	return actions
}

func (g *GullyRestorationGroup) createManagementActions() {
	g.actionMap = make(map[planningunit.Id]*GullyRestoration)
	for planningUnit := range g.sedimentContribution.contributionMap {
		g.createManagementAction(planningUnit)
	}
}

func (g *GullyRestorationGroup) createManagementAction(planningUnit planningunit.Id) {
	originalGullySediment := g.sedimentContribution.SedimentContribution(planningUnit)

	costInDollars := g.implementationCost(planningUnit)
	opportunityCostInDollars := g.opportunityCost(planningUnit)

	actionedGullySedimentReduction := 1 - g.parameters.GetFloat64(parameters.GullySedimentReductionTarget)

	g.actionMap[planningUnit] =
		NewGullyRestoration().
			WithPlanningUnit(planningUnit).
			WithOriginalGullySediment(originalGullySediment).
			WithActionedGullySediment(actionedGullySedimentReduction * originalGullySediment).
			WithImplementationCost(costInDollars).
			WithOpportunityCost(opportunityCostInDollars)
}
