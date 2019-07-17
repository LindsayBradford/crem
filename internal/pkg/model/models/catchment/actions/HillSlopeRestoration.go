// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

const HillSlopeRestorationType action.ManagementActionType = "HillSlopeRestoration"

func NewHillSlopeRestoration() *HillSlopeRestoration {
	action := new(HillSlopeRestoration)
	action.WithType(HillSlopeRestorationType)
	return action
}

type HillSlopeRestoration struct {
	action.SimpleManagementAction
}

func (h *HillSlopeRestoration) WithPlanningUnit(planningUnit planningunit.Id) *HillSlopeRestoration {
	h.SimpleManagementAction.WithPlanningUnit(planningUnit)
	return h
}

const HillSlopeRestorationCost action.ModelVariableName = "HillSlopeRestorationCost"

func (h *HillSlopeRestoration) WithImplementationCost(costInDollars float64) *HillSlopeRestoration {
	h.SimpleManagementAction.WithVariable(HillSlopeRestorationCost, costInDollars)
	return h
}

const OriginalHillSlopeVegetation action.ModelVariableName = "OriginalHillSlopeVegetation"

func (h *HillSlopeRestoration) WithOriginalHillSlopeVegetation(proportionOfVegetation float64) *HillSlopeRestoration {
	h.SimpleManagementAction.WithVariable(OriginalHillSlopeVegetation, proportionOfVegetation)
	return h
}

const ActionedHillSlopeVegetation action.ModelVariableName = "ActionedHillSlopeVegetation"

func (h *HillSlopeRestoration) WithActionedHillSlopeVegetation(proportionOfVegetation float64) *HillSlopeRestoration {
	h.SimpleManagementAction.WithVariable(ActionedHillSlopeVegetation, proportionOfVegetation)
	return h
}
