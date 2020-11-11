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

func (h *HillSlopeRestoration) WithOpportunityCost(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(HillSlopeRestorationOpportunityCost, costInDollars)
}

func (h *HillSlopeRestoration) WithOriginalSedimentErosion(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(HillSlopeErosionOriginalAttribute, costInDollars)
}

func (h *HillSlopeRestoration) WithActionedSedimentErosion(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(HillSlopeErosionActionedAttribute, costInDollars)
}

func (h *HillSlopeRestoration) WithOriginalParticulateNitrogen(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(ParticulateNitrogenOriginalAttribute, costInDollars)
}

func (h *HillSlopeRestoration) WithActionedParticulateNitrogen(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(ParticulateNitrogenActionedAttribute, costInDollars)
}

func (h *HillSlopeRestoration) WithOriginalDissolvedNitrogen(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(DissolvedNitrogenOriginalAttribute, costInDollars)
}

func (h *HillSlopeRestoration) WithActionedDissolvedNitrogen(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(DissolvedNitrogenActionedAttribute, costInDollars)
}

func (h *HillSlopeRestoration) WithVariable(variableName action.ModelVariableName, value float64) *HillSlopeRestoration {
	h.SimpleManagementAction.WithVariable(variableName, value)
	return h
}
