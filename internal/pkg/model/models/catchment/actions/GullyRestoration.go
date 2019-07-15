// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

type GullyRestoration struct {
	action.SimpleManagementAction
}

func (g *GullyRestoration) WithPlanningUnit(planningUnit planningunit.Id) *GullyRestoration {
	g.SimpleManagementAction.WithPlanningUnit(planningUnit)
	return g
}

const GullyRestorationType action.ManagementActionType = "GullyRestoration"

func (g *GullyRestoration) WithGullyRestorationType() *GullyRestoration {
	g.SimpleManagementAction.WithType(GullyRestorationType)
	return g
}

const GullyRestorationCost action.ModelVariableName = "GullyRestorationCost"

func (r *GullyRestoration) WithImplementationCost(costInDollars float64) *GullyRestoration {
	r.SimpleManagementAction.WithVariable(GullyRestorationCost, costInDollars)
	return r
}

const OriginalGullySediment action.ModelVariableName = "OriginalGullySediment"

func (r *GullyRestoration) WithOriginalGullySediment(gullyVolume float64) *GullyRestoration {
	r.SimpleManagementAction.WithVariable(OriginalGullySediment, gullyVolume)
	return r
}

const ActionedGullySediment action.ModelVariableName = "ActionedGullySediment"

func (r *GullyRestoration) WithActionedGullySediment(gullyVolume float64) *GullyRestoration {
	r.SimpleManagementAction.WithVariable(ActionedGullySediment, gullyVolume)
	return r
}
