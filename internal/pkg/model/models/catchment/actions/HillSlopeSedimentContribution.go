package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
)

const (
	hillSlopeAreaIndex                   = 11
	proportionOfHillSlopeVegetationIndex = 12
	HillslopeRKLSIndex                   = 13
)

type hillslopeSedimentTracker struct {
	rslk                         float64
	originalProportionVegetation float64
}

type HillSlopeSedimentContribution struct {
	planningUnitTable tables.CsvTable
	parameters        parameters.Parameters

	contributionMap map[planningunit.Id]hillslopeSedimentTracker
}

func (h *HillSlopeSedimentContribution) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) {
	h.planningUnitTable = planningUnitTable
	h.parameters = parameters
	h.populateContributionMap()
}

func (h *HillSlopeSedimentContribution) populateContributionMap() {
	_, rowCount := h.planningUnitTable.ColumnAndRowSize()
	h.contributionMap = make(map[planningunit.Id]hillslopeSedimentTracker, rowCount)

	for row := uint(0); row < rowCount; row++ {
		h.populateContributionMapEntry(row)
	}
}

func (h *HillSlopeSedimentContribution) populateContributionMapEntry(rowNumber uint) {
	planningUnit := h.planningUnitTable.CellFloat64(planningUnitIndex, rowNumber)
	mapKey := planningunit.Float64ToId(planningUnit)

	h.contributionMap[mapKey] = hillslopeSedimentTracker{
		rslk:                         h.hillslopeRkls(rowNumber),
		originalProportionVegetation: h.originalHillSlopeVegetation(rowNumber),
	}
}

func (h *HillSlopeSedimentContribution) hillslopeRkls(rowNumber uint) float64 {
	// rkls: Rainfall erosivity factor  (R) * Soil Erodibility Factor (K) * Slope length (L) * Slope Steepness (S)
	// See Catchment Rehabilitation Planner final report, section 3.2.3
	rkls := h.planningUnitTable.CellFloat64(HillslopeRKLSIndex, rowNumber)
	return rkls
}

func (h *HillSlopeSedimentContribution) originalHillSlopeVegetation(rowNumber uint) float64 {
	return h.planningUnitTable.CellFloat64(proportionOfHillSlopeVegetationIndex, rowNumber)
}

func (h *HillSlopeSedimentContribution) OriginalSedimentContribution() float64 {
	sedimentContribution := float64(0)
	for planningUnit := range h.contributionMap {
		sedimentContribution += h.OriginalPlanningUnitSedimentContribution(planningUnit)
	}
	return sedimentContribution
}

func (h *HillSlopeSedimentContribution) OriginalPlanningUnitSedimentContribution(id planningunit.Id) float64 {
	// TODO:  Need to reduce if matching RiverBankingRestoration is active.
	planningUnitSedimentTracker, planningUnitIsPresent := h.contributionMap[id]
	assert.That(planningUnitIsPresent).Holds()

	sedimentContribution := planningUnitSedimentTracker.rslk * planningUnitSedimentTracker.originalProportionVegetation

	return sedimentContribution
}

func (h *HillSlopeSedimentContribution) PlanningUnitSedimentContribution(planningUnit planningunit.Id, proportionOfHillSlopeVegetation float64) float64 {
	// TODO:  Need to reduce if matching RiverBankingRestoration is active.
	planningUnitSedimentTracker, planningUnitIsPresent := h.contributionMap[planningUnit]
	assert.That(planningUnitIsPresent).Holds()

	sedimentContribution := planningUnitSedimentTracker.rslk * proportionOfHillSlopeVegetation

	return sedimentContribution
}
