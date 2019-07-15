// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
)

const (
	planningUnitIndex    = 0
	riverLengthIndex     = 2
	channelSlopeIndex    = 3
	bankFullFlowIndex    = 4
	bankHeightIndex      = 6
	floodPlainWidthIndex = 7

	riparianVegetationIndex = 8
	subCatchmentAreaIndex   = 9
	riparianBufferAreaIndex = 10
)

func Float64ToPlanningUnitId(value float64) planningunit.Id {
	return planningunit.Id(value)
}

type sedimentTracker struct {
	partialSedimentContribution      float64
	originalIntactRiparianVegetation float64
}

type BankSedimentContribution struct {
	planningUnitTable tables.CsvTable
	parameters        parameters.Parameters

	contributionMap map[planningunit.Id]sedimentTracker
}

func (bsc *BankSedimentContribution) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) {
	bsc.planningUnitTable = planningUnitTable
	bsc.parameters = parameters
	bsc.populateContributionMap()
}

func (bsc *BankSedimentContribution) populateContributionMap() {
	_, rowCount := bsc.planningUnitTable.ColumnAndRowSize()
	bsc.contributionMap = make(map[planningunit.Id]sedimentTracker, rowCount)

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

	sedimentDensity := bsc.parameters.GetFloat64(parameters.SedimentDensity)
	suspendedSedimentProportion := bsc.parameters.GetFloat64(parameters.SuspendedSedimentProportion)

	bankErosionFudgeFactor := bsc.parameters.GetFloat64(parameters.BankErosionFudgeFactor)

	waterDensity := bsc.parameters.GetFloat64(parameters.WaterDensity)
	localAcceleration := bsc.parameters.GetFloat64(parameters.LocalAcceleration)
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

func (bsc *BankSedimentContribution) OriginalPlanningUnitSedimentContribution(id planningunit.Id) float64 {
	planningUnitSedimentTracker, planningUnitIsPresent := bsc.contributionMap[id]
	assert.That(planningUnitIsPresent).Holds()

	originalVegetation := planningUnitSedimentTracker.originalIntactRiparianVegetation
	sedimentContribution :=
		planningUnitSedimentTracker.partialSedimentContribution * bsc.adjustedProportionOfIntactVegetation(originalVegetation)

	return sedimentContribution
}

func (bsc *BankSedimentContribution) PlanningUnitSedimentContribution(planningUnit planningunit.Id, proportionOfIntactVegetation float64) float64 {
	planningUnitSedimentTracker, planningUnitIsPresent := bsc.contributionMap[planningUnit]
	assert.That(planningUnitIsPresent).Holds()

	sedimentContribution :=
		planningUnitSedimentTracker.partialSedimentContribution * bsc.adjustedProportionOfIntactVegetation(proportionOfIntactVegetation)

	return sedimentContribution
}

func (bsc *BankSedimentContribution) adjustedProportionOfIntactVegetation(proportionOfIntactVegetation float64) float64 {
	return 1 - 0.95*proportionOfIntactVegetation //TODO: Why a 5% dampener here?  Nothing in documentation. Should it be a parameter?
}
