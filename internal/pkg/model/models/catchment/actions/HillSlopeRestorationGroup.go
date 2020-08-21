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

	if !h.actionNeededFor(planningUnitAsId) {
		return
	}

	originalHillSlopeErosion := h.originalHillSlopeErosion(planningUnitAsId)
	adjustedOriginalHillSlopeErosion := h.adjustByDeliveryRatio(originalHillSlopeErosion)
	actionedHillSlopeErosion := h.actionedHillSlopeErosion(planningUnitAsId)
	adjustedActionedHillSlopeErosion := h.adjustByDeliveryRatio(actionedHillSlopeErosion)

	opportunityCostInDollars := h.opportunityCost(planningUnitAsId)
	implementationCostInDollars := h.implementationCost(planningUnitAsId)

	originalParticulateNitrogen := h.originalParticulateNitrogen(planningUnitAsId)
	adjustedOriginalParticulateNitrogen := h.adjustByDeliveryRatio(originalParticulateNitrogen)
	actionedParticulateNitrogen := h.actionedParticulateNitrogen(planningUnitAsId)
	adjustedActionedParticulateNitrogen := h.adjustByDeliveryRatio(actionedParticulateNitrogen)

	originalFineSediment := h.originalFineSediment(planningUnitAsId)
	actionedFineSediment := h.originalFineSediment(planningUnitAsId)

	h.actionMap[planningUnitAsId] =
		NewHillSlopeRestoration().
			WithPlanningUnit(planningUnitAsId).
			WithOriginalSedimentErosion(adjustedOriginalHillSlopeErosion).
			WithActionedSedimentErosion(adjustedActionedHillSlopeErosion).
			WithOriginalParticulateNitrogen(adjustedOriginalParticulateNitrogen).
			WithActionedParticulateNitrogen(adjustedActionedParticulateNitrogen).
			WithOriginalFineSediment(originalFineSediment).
			WithActionedFineSediment(actionedFineSediment).
			WithOpportunityCost(opportunityCostInDollars).
			WithImplementationCost(implementationCostInDollars)
}

func (h *HillSlopeRestorationGroup) actionNeededFor(planningUnit planningunit.Id) bool {
	const minimumPrecision = 3
	originalSediment := h.originalHillSlopeErosion(planningUnit)
	roundedSediment := math.RoundFloat(originalSediment, minimumPrecision)

	return h.originalHillSlopeErosion(planningUnit) > 0 && roundedSediment > 0
}

func (h *HillSlopeRestorationGroup) adjustByDeliveryRatio(value float64) float64 {
	sedimentDeliveryRatio := h.parameters.GetFloat64(parameters.HillSlopeDeliveryRatio)
	return value * sedimentDeliveryRatio
}
