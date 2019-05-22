// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
)

const actionedGullyVolume = 0

type GullyRestorations struct {
	sedimentContribution *GullySedimentContribution
	parameters           parameters.Parameters

	actionMap map[string]*GullyRestoration
}

func (g *GullyRestorations) Initialise(gullyTable tables.CsvTable, parameters parameters.Parameters) *GullyRestorations {
	g.sedimentContribution = new(GullySedimentContribution)
	g.sedimentContribution.Initialise(gullyTable, parameters)
	g.parameters = parameters
	g.createManagementActions()

	return g
}

func (g *GullyRestorations) ManagementActions() map[string]*GullyRestoration {
	return g.actionMap
}

func (g *GullyRestorations) createManagementActions() {
	g.actionMap = make(map[string]*GullyRestoration)
	for planningUnit := range g.sedimentContribution.contributionMap {
		g.createManagementAction(planningUnit)
	}
}

func (g *GullyRestorations) createManagementAction(planningUnit string) {
	originalGullyVolume := g.sedimentContribution.SedimentContribution(planningUnit)
	costInDollars := g.calculateImplementationCost(planningUnit)

	g.actionMap[planningUnit] =
		new(GullyRestoration).
			WithGullyRestorationType().
			WithPlanningUnit(string(planningUnit)).
			WithOriginalGullyVolume(originalGullyVolume).
			WithActionedGullyVolume(actionedGullyVolume).
			WithImplementationCost(costInDollars)
}

func (g *GullyRestorations) calculateImplementationCost(planningUnit string) float64 {
	channelRestorationCostPerKilometer := g.parameters.GetFloat64(parameters.GullyRestorationCostPerKilometer)

	channelLengthInMetres := g.sedimentContribution.ChannelLength(planningUnit)
	channelLengthInKilometres := channelLengthInMetres / 1000

	implementationCost := channelLengthInKilometres * channelRestorationCostPerKilometer

	return implementationCost
}
