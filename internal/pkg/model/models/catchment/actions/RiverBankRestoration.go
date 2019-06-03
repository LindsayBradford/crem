// Copyright (c) 2019 Australian Rivers Institute.

package actions

import "github.com/LindsayBradford/crem/internal/pkg/model/action"

type RiverBankRestoration struct {
	action.SimpleManagementAction
}

func (r *RiverBankRestoration) WithPlanningUnit(planningUnit string) *RiverBankRestoration {
	r.SimpleManagementAction.WithPlanningUnit(planningUnit)
	return r
}

const RiverBankRestorationType action.ManagementActionType = "RiverBankRestoration"

func (r *RiverBankRestoration) WithRiverBankRestorationType() *RiverBankRestoration {
	r.SimpleManagementAction.WithType(RiverBankRestorationType)
	return r
}

const RiverBankRestorationCost action.ModelVariableName = "RiverBankRestorationCost"

func (r *RiverBankRestoration) WithImplementationCost(costInDollars float64) *RiverBankRestoration {
	r.SimpleManagementAction.WithVariable(RiverBankRestorationCost, costInDollars)
	return r
}

const ActionedBufferVegetation action.ModelVariableName = "ActionedBufferVegetation"

func (r *RiverBankRestoration) WithActionedBufferVegetation(changeAsBufferProportion float64) *RiverBankRestoration {
	r.SimpleManagementAction.WithVariable(ActionedBufferVegetation, changeAsBufferProportion)
	return r
}

const OriginalBufferVegetation action.ModelVariableName = "OriginalBufferVegetation"

func (r *RiverBankRestoration) WithUnActionedBufferVegetation(originalBufferProportion float64) *RiverBankRestoration {
	r.SimpleManagementAction.WithVariable(OriginalBufferVegetation, originalBufferProportion)
	return r
}
