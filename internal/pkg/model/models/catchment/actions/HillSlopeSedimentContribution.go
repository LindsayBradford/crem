package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
)

const (
	hillSlopeAreaIndex = 11
)

type hillSlopeSedimentTracker struct {
	area                     float64
	originalSedimentProduced float64
	actionedSedimentProduced float64
}

type HillSlopeSedimentContribution struct {
	planningUnitTable tables.CsvTable
	parameters        parameters.Parameters

	contributionMap       map[planningunit.Id]hillSlopeSedimentTracker
	sedimentDeliveryRatio float64
}

func (h *HillSlopeSedimentContribution) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) {
	h.planningUnitTable = planningUnitTable
	h.parameters = parameters
	h.populateContributionMap()
}

func (h *HillSlopeSedimentContribution) populateContributionMap() {
	h.sedimentDeliveryRatio = h.parameters.GetFloat64(parameters.HillSlopeDeliveryRatio)
	_, rowCount := h.planningUnitTable.ColumnAndRowSize()
	h.contributionMap = make(map[planningunit.Id]hillSlopeSedimentTracker, rowCount)

	for row := uint(0); row < rowCount; row++ {
		h.populateContributionMapEntry(row)
	}
}

func (h *HillSlopeSedimentContribution) populateContributionMapEntry(rowNumber uint) {
	planningUnit := h.planningUnitTable.CellFloat64(planningUnitIndex, rowNumber)
	mapKey := planningunit.Float64ToId(planningUnit)

	h.contributionMap[mapKey] = hillSlopeSedimentTracker{
		area:                     h.hillSlopeArea(rowNumber),
		originalSedimentProduced: 1,
		actionedSedimentProduced: 0,
	}
}

func (h *HillSlopeSedimentContribution) hillSlopeArea(rowNumber uint) float64 {
	return h.planningUnitTable.CellFloat64(hillSlopeAreaIndex, rowNumber)
}

func (h *HillSlopeSedimentContribution) OriginalPlanningUnitSedimentContribution(id planningunit.Id) float64 {
	sedimentTracker, planningUnitIsPresent := h.contributionMap[id]
	assert.That(planningUnitIsPresent).Holds()

	originalSediment := h.calculateDeliveryAdjustedSediment(sedimentTracker.originalSedimentProduced)
	return originalSediment
}

func (h *HillSlopeSedimentContribution) PlanningUnitSedimentContribution(id planningunit.Id, rawSedimentProduced float64) float64 {
	_, planningUnitIsPresent := h.contributionMap[id]
	assert.That(planningUnitIsPresent).Holds()

	originalSediment := h.calculateDeliveryAdjustedSediment(rawSedimentProduced)
	return originalSediment
}

func (h *HillSlopeSedimentContribution) calculateDeliveryAdjustedSediment(sedimentProduced float64) float64 {
	return sedimentProduced * h.sedimentDeliveryRatio
}
