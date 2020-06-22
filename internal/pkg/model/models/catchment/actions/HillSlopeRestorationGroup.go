// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

type HillSlopeRestorationGroup struct {
	planningUnitTable tables.CsvTable
	parameters        parameters.Parameters

	actionMap map[planningunit.Id]*HillSlopeRestoration
	Container
}

func (h *HillSlopeRestorationGroup) WithParameters(parameters parameters.Parameters) *HillSlopeRestorationGroup {
	h.parameters = parameters
	return h
}

func (h *HillSlopeRestorationGroup) WithPlanningUnitTable(planningUnitTable tables.CsvTable) *HillSlopeRestorationGroup {
	h.planningUnitTable = planningUnitTable
	return h
}

func (h *HillSlopeRestorationGroup) WithParentSoilsTable(parentSoilsTable tables.CsvTable) *HillSlopeRestorationGroup {
	h.Container.WithSourceFilter(HillSlopeSource).WithActionsTable(parentSoilsTable)
	return h
}

func (h *HillSlopeRestorationGroup) ManagementActions() []action.ManagementAction {
	h.createManagementActions()
	actions := make([]action.ManagementAction, 0)
	for _, value := range h.actionMap {
		actions = append(actions, value)
	}
	return actions
}

func (h *HillSlopeRestorationGroup) createManagementActions() {
	_, rowCount := h.planningUnitTable.ColumnAndRowSize()
	h.actionMap = make(map[planningunit.Id]*HillSlopeRestoration, rowCount)

	for row := uint(0); row < rowCount; row++ {
		h.createManagementAction(row)
	}
}

func (h *HillSlopeRestorationGroup) createManagementAction(rowNumber uint) {
	planningUnit := h.planningUnitTable.CellFloat64(planningUnitIndex, rowNumber)
	planningUnitAsId := planningunit.Float64ToId(planningUnit)

	hillSlopeArea := h.planningUnitTable.CellFloat64(hillSlopeAreaIndex, rowNumber)
	vegetationTarget := h.parameters.GetFloat64(parameters.HillSlopeBevegetationProportionTarget)
	originalHillSlopeVegetation := h.originalHillSlopeVegetation(rowNumber)

	if hillSlopeArea == 0 || originalHillSlopeVegetation >= vegetationTarget {
		return
	}

	nitrogenValue := h.nitrogenAttributeValue(planningUnitAsId)
	carbonValue := h.carbonAttributeValue(planningUnitAsId)
	deltaCarbonValue := h.deltaCarbonAttributeValue(planningUnitAsId)
	actionedCarbonValue := carbonValue + deltaCarbonValue

	costInDollars := h.calculateImplementationCost(rowNumber)
	opportunityCostInDollars := h.calculateImplementationCost(rowNumber)

	h.actionMap[planningUnitAsId] =
		NewHillSlopeRestoration().
			WithPlanningUnit(planningUnitAsId).
			WithOriginalHillSlopeVegetation(originalHillSlopeVegetation).
			WithActionedHillSlopeVegetation(vegetationTarget).
			WithTotalNitrogen(nitrogenValue).
			WithOriginalTotalCarbon(carbonValue).
			WithActionedTotalCarbon(actionedCarbonValue).
			WithImplementationCost(costInDollars).
			WithOpportunityCost(opportunityCostInDollars)
}

func (h *HillSlopeRestorationGroup) originalHillSlopeVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := h.planningUnitTable.CellFloat64(proportionOfHillSlopeVegetationIndex, rowNumber)
	return proportionOfRiparianVegetation
}

func (h *HillSlopeRestorationGroup) calculateChangeInHillSlopeVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := h.originalHillSlopeVegetation(rowNumber)
	vegetationTarget := h.parameters.GetFloat64(parameters.HillSlopeBevegetationProportionTarget)
	changeInRiparianVegetation := vegetationTarget - proportionOfRiparianVegetation
	return changeInRiparianVegetation
}

func (h *HillSlopeRestorationGroup) calculateImplementationCost(rowNumber uint) float64 {
	implementationCostPerKmSquared := h.parameters.GetFloat64(parameters.HillSlopeRestorationCostPerKilometerSquared)
	hillSlopeAreaInMetresSquared := h.planningUnitTable.CellFloat64(hillSlopeAreaIndex, rowNumber)
	hillSlopeAreaInKilometresSquared := hillSlopeAreaInMetresSquared / 1000

	vegetationChange := h.calculateChangeInHillSlopeVegetation(rowNumber)

	vegetationChangeInKilometresSquared := vegetationChange * hillSlopeAreaInKilometresSquared

	implementationCost := vegetationChangeInKilometresSquared * implementationCostPerKmSquared

	return implementationCost
}
