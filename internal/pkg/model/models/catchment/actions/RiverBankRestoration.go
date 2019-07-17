// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

const RiverBankRestorationType action.ManagementActionType = "RiverBankRestoration"

func NewRiverBankRestoration() *RiverBankRestoration {
	action := new(RiverBankRestoration)
	action.WithType(RiverBankRestorationType)
	return action
}

type RiverBankRestoration struct {
	action.SimpleManagementAction
}

func (r *RiverBankRestoration) WithPlanningUnit(planningUnit planningunit.Id) *RiverBankRestoration {
	r.SimpleManagementAction.WithPlanningUnit(planningUnit)
	return r
}

const RiverBankRestorationCost action.ModelVariableName = "RiverBankRestorationCost"

func (r *RiverBankRestoration) WithImplementationCost(costInDollars float64) *RiverBankRestoration {
	r.SimpleManagementAction.WithVariable(RiverBankRestorationCost, costInDollars)
	return r
}

const ActionedBufferVegetation action.ModelVariableName = "ActionedBufferVegetation"

func (r *RiverBankRestoration) WithActionedBufferVegetation(proportionOfVegetation float64) *RiverBankRestoration {
	r.SimpleManagementAction.WithVariable(ActionedBufferVegetation, proportionOfVegetation)
	return r
}

const OriginalBufferVegetation action.ModelVariableName = "OriginalBufferVegetation"

func (r *RiverBankRestoration) WithUnActionedBufferVegetation(proportionOfVegetation float64) *RiverBankRestoration {
	r.SimpleManagementAction.WithVariable(OriginalBufferVegetation, proportionOfVegetation)
	return r
}
