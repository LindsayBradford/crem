// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

const GullyRestorationType action.ManagementActionType = "GullyRestoration"

func NewGullyRestoration() *GullyRestoration {
	return new(GullyRestoration).WithType(GullyRestorationType)
}

type GullyRestoration struct {
	action.SimpleManagementAction
}

func (g *GullyRestoration) WithType(actionType action.ManagementActionType) *GullyRestoration {
	g.SimpleManagementAction.WithType(actionType)
	return g
}

func (g *GullyRestoration) WithPlanningUnit(planningUnit planningunit.Id) *GullyRestoration {
	g.SimpleManagementAction.WithPlanningUnit(planningUnit)
	return g
}

const GullyRestorationCost action.ModelVariableName = "GullyRestorationCost"

func (g *GullyRestoration) WithImplementationCost(costInDollars float64) *GullyRestoration {
	return g.WithVariable(GullyRestorationCost, costInDollars)
}

const GullyRestorationOpportunityCost action.ModelVariableName = "GullyRestorationOpportunityCost"

func (g *GullyRestoration) WithOpportunityCost(costInDollars float64) *GullyRestoration {
	return g.WithVariable(GullyRestorationOpportunityCost, costInDollars)
}

const OriginalGullySediment action.ModelVariableName = "OriginalGullySediment"

func (g *GullyRestoration) WithOriginalGullySediment(gullyVolume float64) *GullyRestoration {
	return g.WithVariable(OriginalGullySediment, gullyVolume)
}

const ActionedGullySediment action.ModelVariableName = "ActionedGullySediment"

func (g *GullyRestoration) WithActionedGullySediment(gullyVolume float64) *GullyRestoration {
	return g.WithVariable(ActionedGullySediment, gullyVolume)
}

func (g *GullyRestoration) WithVariable(variableName action.ModelVariableName, value float64) *GullyRestoration {
	g.SimpleManagementAction.WithVariable(variableName, value)
	return g
}
