// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

const HillSlopeRestorationType action.ManagementActionType = "HillSlopeRestoration"

func NewHillSlopeRestoration() *HillSlopeRestoration {
	return new(HillSlopeRestoration).WithType(HillSlopeRestorationType)
}

type HillSlopeRestoration struct {
	action.SimpleManagementAction
}

func (h *HillSlopeRestoration) WithType(actionType action.ManagementActionType) *HillSlopeRestoration {
	h.SimpleManagementAction.WithType(actionType)
	return h
}

func (h *HillSlopeRestoration) WithPlanningUnit(planningUnit planningunit.Id) *HillSlopeRestoration {
	h.SimpleManagementAction.WithPlanningUnit(planningUnit)
	return h
}

const HillSlopeRestorationCost action.ModelVariableName = "HillSlopeRestorationCost"

func (h *HillSlopeRestoration) WithImplementationCost(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(HillSlopeRestorationCost, costInDollars)
}

const HillSlopeRestorationOpportunityCost action.ModelVariableName = "HillSlopeRestorationOpportunityCost"

func (g *HillSlopeRestoration) WithOpportunityCost(costInDollars float64) *HillSlopeRestoration {
	return g.WithVariable(HillSlopeRestorationOpportunityCost, costInDollars)
}

const OriginalHillSlopeVegetation action.ModelVariableName = "OriginalHillSlopeVegetation"

func (h *HillSlopeRestoration) WithOriginalHillSlopeVegetation(proportionOfVegetation float64) *HillSlopeRestoration {
	return h.WithVariable(OriginalHillSlopeVegetation, proportionOfVegetation)
}

const ActionedHillSlopeVegetation action.ModelVariableName = "ActionedHillSlopeVegetation"

func (h *HillSlopeRestoration) WithActionedHillSlopeVegetation(proportionOfVegetation float64) *HillSlopeRestoration {
	return h.WithVariable(ActionedHillSlopeVegetation, proportionOfVegetation)
}

func (h *HillSlopeRestoration) WithTotalNitrogen(totalNitrogen float64) *HillSlopeRestoration {
	return h.WithVariable(TotalNitrogen, totalNitrogen)
}

func (h *HillSlopeRestoration) WithOriginalTotalCarbon(totalCarbon float64) *HillSlopeRestoration {
	return h.WithVariable(OriginalTotalCarbon, totalCarbon)
}

func (h *HillSlopeRestoration) WithActionedTotalCarbon(deltaCarbon float64) *HillSlopeRestoration {
	return h.WithVariable(ActionedTotalCarbon, deltaCarbon)
}

func (h *HillSlopeRestoration) WithVariable(variableName action.ModelVariableName, value float64) *HillSlopeRestoration {
	h.SimpleManagementAction.WithVariable(variableName, value)
	return h
}
