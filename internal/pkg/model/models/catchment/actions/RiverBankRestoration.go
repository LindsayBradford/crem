// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

const RiverBankRestorationType action.ManagementActionType = "RiverBankRestoration"

func NewRiverBankRestoration() *RiverBankRestoration {
	return new(RiverBankRestoration).WithType(RiverBankRestorationType)
}

type RiverBankRestoration struct {
	action.SimpleManagementAction
}

func (r *RiverBankRestoration) WithType(actionType action.ManagementActionType) *RiverBankRestoration {
	r.SimpleManagementAction.WithType(actionType)
	return r
}

func (r *RiverBankRestoration) WithPlanningUnit(planningUnit planningunit.Id) *RiverBankRestoration {
	r.SimpleManagementAction.WithPlanningUnit(planningUnit)
	return r
}

const RiverBankRestorationCost action.ModelVariableName = "RiverBankRestorationCost"

func (r *RiverBankRestoration) WithImplementationCost(costInDollars float64) *RiverBankRestoration {
	return r.WithVariable(RiverBankRestorationCost, costInDollars)
}

const RiverBankRestorationOpportunityCost action.ModelVariableName = "RiverBankRestorationOpportunityCost"

func (r *RiverBankRestoration) WithOpportunityCost(costInDollars float64) *RiverBankRestoration {
	return r.WithVariable(RiverBankRestorationOpportunityCost, costInDollars)
}

const ActionedBufferVegetation action.ModelVariableName = "ActionedBufferVegetation"

func (r *RiverBankRestoration) WithActionedBufferVegetation(proportionOfVegetation float64) *RiverBankRestoration {
	return r.WithVariable(ActionedBufferVegetation, proportionOfVegetation)
}

const OriginalBufferVegetation action.ModelVariableName = "OriginalBufferVegetation"

func (r *RiverBankRestoration) WithOriginalBufferVegetation(proportionOfVegetation float64) *RiverBankRestoration {
	return r.WithVariable(OriginalBufferVegetation, proportionOfVegetation)
}

const OriginalRiparianSedimentProduction action.ModelVariableName = "OriginalRiparianSedimentProduction"

func (r *RiverBankRestoration) WithOriginalRiparianSedimentProduction(sediment float64) *RiverBankRestoration {
	return r.WithVariable(OriginalRiparianSedimentProduction, sediment)
}

const ActionedRiparianSedimentProduction action.ModelVariableName = "ActionedRiparianSedimentProduction"

func (r *RiverBankRestoration) WithActionedRiparianSedimentProduction(sediment float64) *RiverBankRestoration {
	return r.WithVariable(ActionedRiparianSedimentProduction, sediment)
}

func (r *RiverBankRestoration) WithOriginalParticulateNitrogen(particulateNitrogen float64) *RiverBankRestoration {
	return r.WithVariable(ParticulateNitrogenOriginalAttribute, particulateNitrogen)
}

func (r *RiverBankRestoration) WithActionedParticulateNitrogen(particulateNitrogen float64) *RiverBankRestoration {
	return r.WithVariable(ParticulateNitrogenActionedAttribute, particulateNitrogen)
}

func (r *RiverBankRestoration) WithOriginalFineSediment(fineSediment float64) *RiverBankRestoration {
	return r.WithVariable(FineSedimentOriginalAttribute, fineSediment)
}

func (r *RiverBankRestoration) WithActionedFineSediment(fineSediment float64) *RiverBankRestoration {
	return r.WithVariable(FineSedimentActionedAttribute, fineSediment)
}

func (r *RiverBankRestoration) WithVariable(variableName action.ModelVariableName, value float64) *RiverBankRestoration {
	r.SimpleManagementAction.WithVariable(variableName, value)
	return r
}
