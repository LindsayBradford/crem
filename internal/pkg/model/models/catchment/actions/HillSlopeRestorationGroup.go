// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/pkg/math"
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

func (h *HillSlopeRestorationGroup) WithActionsTable(actionsTable tables.CsvTable) *HillSlopeRestorationGroup {
	h.Container.WithSourceFilter(HillSlopeSource).WithActionsTable(actionsTable)
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

	originalBufferVegetation := h.originalBufferVegetation(rowNumber)
	riparianFilter := riparianBufferFilter(originalBufferVegetation)

	if !h.actionNeededFor(planningUnitAsId, riparianFilter) {
		return
	}

	hillSlopeDeliveryRatio := h.parameters.GetFloat64(parameters.HillSlopeDeliveryRatio)

	originalHillSlopeErosion := h.originalHillSlopeErosion(planningUnitAsId) * hillSlopeDeliveryRatio
	actionedHillSlopeErosion := h.actionedHillSlopeErosion(planningUnitAsId) * hillSlopeDeliveryRatio

	opportunityCostInDollars := h.opportunityCost(planningUnitAsId)
	implementationCostInDollars := h.implementationCost(planningUnitAsId)

	originalParticulateNitrogen := h.originalParticulateNitrogen(planningUnitAsId) * hillSlopeDeliveryRatio
	actionedParticulateNitrogen := h.actionedParticulateNitrogen(planningUnitAsId) * hillSlopeDeliveryRatio

	h.actionMap[planningUnitAsId] =
		NewHillSlopeRestoration().
			WithPlanningUnit(planningUnitAsId).
			WithOriginalSedimentErosion(originalHillSlopeErosion).
			WithActionedSedimentErosion(actionedHillSlopeErosion).
			WithOriginalParticulateNitrogen(originalParticulateNitrogen).
			WithActionedParticulateNitrogen(actionedParticulateNitrogen).
			WithOpportunityCost(opportunityCostInDollars).
			WithImplementationCost(implementationCostInDollars)
}

func (h *HillSlopeRestorationGroup) actionNeededFor(planningUnit planningunit.Id, worstCaseRiparianFilter float64) bool {
	originalHillSlopeSediment := h.originalHillSlopeErosion(planningUnit)
	if originalHillSlopeSediment == 0 {
		return false
	}

	const minimumPrecision = 3

	worstCaseHillSlopeSediment := originalHillSlopeSediment *
		h.parameters.GetFloat64(parameters.HillSlopeDeliveryRatio) * worstCaseRiparianFilter

	roundedWorstCaseSediment := math.RoundFloat(worstCaseHillSlopeSediment, minimumPrecision)

	return roundedWorstCaseSediment > 0
}

func (h *HillSlopeRestorationGroup) originalBufferVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := h.planningUnitTable.CellFloat64(riparianVegetationIndex, rowNumber)
	return proportionOfRiparianVegetation
}

func riparianBufferFilter(proportionOfRiparianBufferVegetation float64) float64 {
	if proportionOfRiparianBufferVegetation < 0.25 {
		return 1
	}
	if proportionOfRiparianBufferVegetation > 0.75 {
		return 0.25
	}
	return 1 - proportionOfRiparianBufferVegetation
}
