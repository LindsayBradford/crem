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
}

func (g *GullyRestorationGroup) Initialise(gullyTable tables.CsvTable, parameters parameters.Parameters) *GullyRestorationGroup {
	g.sedimentContribution = new(GullySedimentContribution)
	g.sedimentContribution.Initialise(gullyTable, parameters)
	g.parameters = parameters
	g.createManagementActions()

	return g
}

func (g *GullyRestorationGroup) ManagementActions() []action.ManagementAction {
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
	costInDollars := g.calculateImplementationCost(planningUnit)

	actionedGullySedimentReduction := 1 - g.parameters.GetFloat64(parameters.GullySedimentReductionTarget)

	g.actionMap[planningUnit] =
		NewGullyRestoration().
			WithPlanningUnit(planningUnit).
			WithOriginalGullySediment(originalGullySediment).
			WithActionedGullySediment(actionedGullySedimentReduction * originalGullySediment).
			WithImplementationCost(costInDollars)
}

func (g *GullyRestorationGroup) calculateImplementationCost(planningUnit planningunit.Id) float64 {
	channelRestorationCostPerKilometer := g.parameters.GetFloat64(parameters.GullyRestorationCostPerKilometer)

	channelLengthInMetres := g.sedimentContribution.ChannelLength(planningUnit)
	channelLengthInKilometres := channelLengthInMetres / 1000

	implementationCost := channelLengthInKilometres * channelRestorationCostPerKilometer

	return implementationCost
}
