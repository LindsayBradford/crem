// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

type RiverBankRestorationGroup struct {
	planningUnitTable        tables.CsvTable
	parameters               parameters.Parameters
	bankSedimentContribution BankSedimentContribution

	actionMap map[planningunit.Id]*RiverBankRestoration
	Container
}

func (r *RiverBankRestorationGroup) WithPlanningUnitTable(planningUnitTable tables.CsvTable) *RiverBankRestorationGroup {
	r.planningUnitTable = planningUnitTable
	return r
}

func (r *RiverBankRestorationGroup) WithActionsTable(parentSoilsTable tables.CsvTable) *RiverBankRestorationGroup {
	r.Container.WithSourceFilter(RiparianSource).WithActionsTable(parentSoilsTable)
	return r
}

func (r *RiverBankRestorationGroup) WithParameters(parameters parameters.Parameters) *RiverBankRestorationGroup {
	r.parameters = parameters
	return r
}

func (r *RiverBankRestorationGroup) ManagementActions() []action.ManagementAction {
	r.createManagementActions()
	actions := make([]action.ManagementAction, 0)
	for _, value := range r.actionMap {
		actions = append(actions, value)
	}
	return actions
}

func (r *RiverBankRestorationGroup) createManagementActions() {
	r.bankSedimentContribution.Initialise(r.planningUnitTable, r.parameters)

	_, rowCount := r.planningUnitTable.ColumnAndRowSize()
	r.actionMap = make(map[planningunit.Id]*RiverBankRestoration, rowCount)

	for row := uint(0); row < rowCount; row++ {
		r.createManagementAction(row)
	}
}

func (r *RiverBankRestorationGroup) createManagementAction(rowNumber uint) {
	planningUnit := r.planningUnitTable.CellFloat64(planningUnitIndex, rowNumber)
	planningUnitAsId := planningunit.Float64ToId(planningUnit)

	originalBufferVegetation := r.originalBufferVegetation(rowNumber)
	actionedBufferVegetation := r.parameters.GetFloat64(parameters.RiparianBufferVegetationProportionTarget)

	if originalBufferVegetation >= actionedBufferVegetation {
		return
	}

	opportunityCostInDollars := r.opportunityCost(planningUnitAsId)
	implementationCostInDollars := r.implementationCost(planningUnitAsId)

	originalSediment := r.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnitAsId, originalBufferVegetation)
	actionedSediment := r.bankSedimentContribution.PlanningUnitSedimentContribution(planningUnitAsId, actionedBufferVegetation)

	originalParticulateNitrogen := r.originalParticulateNitrogen(planningUnitAsId)
	actionedParticulateNitrogen := r.actionedParticulateNitrogen(planningUnitAsId)

	originalFineSediment := r.originalFineSediment(planningUnitAsId)
	actionedFineSediment := r.actionedFineSediment(planningUnitAsId)

	r.actionMap[planningUnitAsId] =
		NewRiverBankRestoration().
			WithPlanningUnit(planningUnitAsId).
			WithOriginalBufferVegetation(originalBufferVegetation).
			WithActionedBufferVegetation(actionedBufferVegetation).
			WithOriginalRiparianSedimentProduction(originalSediment).
			WithActionedRiparianSedimentProduction(actionedSediment).
			WithOriginalParticulateNitrogen(originalParticulateNitrogen).
			WithActionedParticulateNitrogen(actionedParticulateNitrogen).
			WithOriginalFineSediment(originalFineSediment).
			WithActionedFineSediment(actionedFineSediment).
			WithImplementationCost(implementationCostInDollars).
			WithOpportunityCost(opportunityCostInDollars)
}

func (r *RiverBankRestorationGroup) originalBufferVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := r.planningUnitTable.CellFloat64(riparianVegetationIndex, rowNumber)
	return proportionOfRiparianVegetation
}

func (r *RiverBankRestorationGroup) calculateChangeInBufferVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := r.originalBufferVegetation(rowNumber)
	vegetationTarget := r.parameters.GetFloat64(parameters.RiparianBufferVegetationProportionTarget)
	changeInRiparianVegetation := vegetationTarget - proportionOfRiparianVegetation
	return changeInRiparianVegetation
}
