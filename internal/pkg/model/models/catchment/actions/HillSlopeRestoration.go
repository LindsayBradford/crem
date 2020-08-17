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

const HillSlopeRestorationOriginalSedimentErosion action.ModelVariableName = "HillSlopeRestorationOriginalSedimentErosion"

func (h *HillSlopeRestoration) WithOriginalSedimentErosion(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(HillSlopeRestorationOriginalSedimentErosion, costInDollars)
}

const HillSlopeRestorationActionedSedimentErosion action.ModelVariableName = "HillSlopeRestorationActionedSedimentErosion"

func (h *HillSlopeRestoration) WithActionedSedimentErosion(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(HillSlopeRestorationActionedSedimentErosion, costInDollars)
}

const HillSlopeRestorationOriginalParticulateNitrogen action.ModelVariableName = "HillSlopeRestorationOriginalParticulateNitrogen"

func (h *HillSlopeRestoration) WithOriginalParticulateNitrogen(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(HillSlopeRestorationOriginalParticulateNitrogen, costInDollars)
}

const HillSlopeRestorationActionedParticulateNitrogen action.ModelVariableName = "HillSlopeRestorationActionedParticulateNitrogen"

func (h *HillSlopeRestoration) WithActionedParticulateNitrogen(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(HillSlopeRestorationActionedParticulateNitrogen, costInDollars)
}

const HillSlopeRestorationOriginalFineSediment action.ModelVariableName = "HillSlopeRestorationOriginalFineSediment"

func (h *HillSlopeRestoration) WithOriginalFineSediment(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(HillSlopeRestorationOriginalFineSediment, costInDollars)
}

const HillSlopeRestorationActionedFineSediment action.ModelVariableName = "HillSlopeRestorationActionedFineSediment"

func (h *HillSlopeRestoration) WithActionedFineSediment(costInDollars float64) *HillSlopeRestoration {
	return h.WithVariable(HillSlopeRestorationActionedFineSediment, costInDollars)
}

func (h *HillSlopeRestoration) WithVariable(variableName action.ModelVariableName, value float64) *HillSlopeRestoration {
	h.SimpleManagementAction.WithVariable(variableName, value)
	return h
}
