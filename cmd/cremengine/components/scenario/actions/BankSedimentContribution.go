// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"math"

	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
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

type planningUnitId int64

type sedimentTracker struct {
	partialSedimentContribution      float64
	originalIntactRiparianVegetation float64
}

type BankSedimentContribution struct {
	planningUnitTable *tables.CsvTable
	parameters        parameters.Parameters

	contributionMap map[planningUnitId]sedimentTracker
}

func (bsc *BankSedimentContribution) Initialise(planningUnitTable *tables.CsvTable, parameters parameters.Parameters) {
	bsc.planningUnitTable = planningUnitTable
	bsc.parameters = parameters
	bsc.populateContributionMap()
}

func (bsc *BankSedimentContribution) populateContributionMap() {
	_, rowCount := bsc.planningUnitTable.Size()
	bsc.contributionMap = make(map[planningUnitId]sedimentTracker, rowCount)

	for row := uint(0); row < rowCount; row++ {
		bsc.populateContributionMapEntry(row)
	}
}

func (bsc *BankSedimentContribution) populateContributionMapEntry(rowNumber uint) {
	planningUnit := bsc.planningUnitTable.CellInt64(planningUnitIndex, rowNumber)
	mapKey := planningUnitId(planningUnit)

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
	for _, tracker := range bsc.contributionMap {
		originalVegetation := tracker.originalIntactRiparianVegetation
		sedimentContribution += tracker.partialSedimentContribution *
			bsc.SedimentImpactOfRiparianVegetation(originalVegetation)
	}
	return sedimentContribution
}

func (bsc *BankSedimentContribution) SedimentImpactOfRiparianVegetation(proportionOfIntactVegetation float64) float64 {
	return 1 - 0.95*proportionOfIntactVegetation //TODO: Why a 5% dampener here?  Nothing in documentation. Should it be a parameter?
}
