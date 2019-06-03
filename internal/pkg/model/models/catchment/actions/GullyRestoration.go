// Copyright (c) 2019 Australian Rivers Institute.

package actions

import "github.com/LindsayBradford/crem/internal/pkg/model/action"

type GullyRestoration struct {
	action.SimpleManagementAction
}

func (g *GullyRestoration) WithPlanningUnit(planningUnit string) *GullyRestoration {
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

const OriginalGullyVolume action.ModelVariableName = "OriginalGullyVolume"

func (r *GullyRestoration) WithOriginalGullyVolume(gullyVolume float64) *GullyRestoration {
	r.SimpleManagementAction.WithVariable(OriginalGullyVolume, gullyVolume)
	return r
}

const ActionedGullyVolume action.ModelVariableName = "ActionedGullyVolume"

func (r *GullyRestoration) WithActionedGullyVolume(gullyVolume float64) *GullyRestoration {
	r.SimpleManagementAction.WithVariable(ActionedGullyVolume, gullyVolume)
	return r
}
