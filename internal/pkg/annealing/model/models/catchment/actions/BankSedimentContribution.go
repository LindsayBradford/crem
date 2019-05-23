// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"math"
	"strconv"

	. "github.com/LindsayBradford/crem/internal/pkg/annealing/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
)

const (
	planningUnitIndex       = 0
	riverLengthIndex        = 2
	channelSlopeIndex       = 3
	bankHeightIndex         = 4
	floodPlainWidthIndex    = 6
	bankFullFlowIndex       = 7
	riparianVegetationIndex = 8
	planningUnitAreaIndex   = 9
	riparianBufferAreaIndex = 10
)

func Float64ToPlanningUnitId(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}

type sedimentTracker struct {
	partialSedimentContribution      float64
	originalIntactRiparianVegetation float64
}

type BankSedimentContribution struct {
	planningUnitTable tables.CsvTable
	parameters        Parameters

	contributionMap map[string]sedimentTracker
}

func (bsc *BankSedimentContribution) Initialise(planningUnitTable tables.CsvTable, parameters Parameters) {
	bsc.planningUnitTable = planningUnitTable
	bsc.parameters = parameters
	bsc.populateContributionMap()
}

func (bsc *BankSedimentContribution) populateContributionMap() {
	_, rowCount := bsc.planningUnitTable.ColumnAndRowSize()
	bsc.contributionMap = make(map[string]sedimentTracker, rowCount)

	for row := uint(0); row < rowCount; row++ {
		bsc.populateContributionMapEntry(row)
	}
}

func (bsc *BankSedimentContribution) populateContributionMapEntry(rowNumber uint) {
	planningUnit := bsc.planningUnitTable.CellFloat64(planningUnitIndex, rowNumber)
	mapKey := Float64ToPlanningUnitId(planningUnit)

	bsc.contributionMap[mapKey] = sedimentTracker{
		partialSedimentContribution:      bsc.partialBankSedimentContribution(rowNumber),
		originalIntactRiparianVegetation: bsc.originalIntactRiparianVegetation(rowNumber),
	}
}

func (bsc *BankSedimentContribution) partialBankSedimentContribution(rowNumber uint) float64 {
	riverLength := bsc.planningUnitTable.CellFloat64(riverLengthIndex, rowNumber)
	bankHeight := bsc.planningUnitTable.CellFloat64(bankHeightIndex, rowNumber)

	sedimentDensity := bsc.parameters.GetFloat64(SedimentDensity)
	suspendedSedimentProportion := bsc.parameters.GetFloat64(SuspendedSedimentProportion)

	bankErosionFudgeFactor := bsc.parameters.GetFloat64(BankErosionFudgeFactor)

	waterDensity := bsc.parameters.GetFloat64(WaterDensity)
	localAcceleration := bsc.parameters.GetFloat64(LocalAcceleration)
	bankFullDischarge := bsc.planningUnitTable.CellFloat64(bankFullFlowIndex, rowNumber)
	channelSlope := bsc.planningUnitTable.CellFloat64(channelSlopeIndex, rowNumber)

	channelDischarge := waterDensity * localAcceleration * bankFullDischarge * channelSlope

	riparianVegetationImpact := float64(1) // This is the value that changes as we anneal, leaving in formula for now for traceability.

	floodPlainWidth := bsc.planningUnitTable.CellFloat64(floodPlainWidthIndex, rowNumber)
	floodPlainWidthRelationship := 1 - math.Exp(-1.5*math.Pow(10, -2.0)*floodPlainWidth)

	return bankErosionFudgeFactor * channelDischarge * riparianVegetationImpact *
		floodPlainWidthRelationship * riverLength * bankHeight * sedimentDensity * suspendedSedimentProportion
}

func (bsc *BankSedimentContribution) originalIntactRiparianVegetation(rowNumber uint) float64 {
	return bsc.planningUnitTable.CellFloat64(riparianVegetationIndex, rowNumber)
}

func (bsc *BankSedimentContribution) OriginalSedimentContribution() float64 {
	sedimentContribution := float64(0)
	for planningUnit := range bsc.contributionMap {
		sedimentContribution += bsc.OriginalPlanningUnitSedimentContribution(planningUnit)
	}
	return sedimentContribution
}

func (bsc *BankSedimentContribution) OriginalPlanningUnitSedimentContribution(id string) float64 {
	planningUnitSedimentTracker, planningUnitIsPresent := bsc.contributionMap[id]
	assert.That(planningUnitIsPresent).Holds()

	originalVegetation := planningUnitSedimentTracker.originalIntactRiparianVegetation
	sedimentContribution :=
		planningUnitSedimentTracker.partialSedimentContribution * bsc.adjustedProportionOfIntactVegetation(originalVegetation)

	return sedimentContribution
}

func (bsc *BankSedimentContribution) PlanningUnitSedimentContribution(planningUnit string, proportionOfIntactVegetation float64) float64 {
	planningUnitSedimentTracker, planningUnitIsPresent := bsc.contributionMap[planningUnit]
	assert.That(planningUnitIsPresent).Holds()

	sedimentContribution :=
		planningUnitSedimentTracker.partialSedimentContribution * bsc.adjustedProportionOfIntactVegetation(proportionOfIntactVegetation)

	return sedimentContribution
}

func (bsc *BankSedimentContribution) adjustedProportionOfIntactVegetation(proportionOfIntactVegetation float64) float64 {
	return 1 - 0.95*proportionOfIntactVegetation //TODO: Why a 5% dampener here?  Nothing in documentation. Should it be a parameter?
}
