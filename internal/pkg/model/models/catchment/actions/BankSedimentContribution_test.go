// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"math"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	. "github.com/onsi/gomega"
)

const expectedRiparianVegetationProportion = float64(0.5)
const equalTo = "=="

const expectedRowNumber = 3
const expectedColumnNumber = 11

const (
	defaultRiverLength        = float64(5)
	defaultRiverSlope         = 0.0025
	defaultBankHeight         = float64(3)
	defaultFloodPlainWidth    = float64(500)
	defaultBankFullFlow       = float64(200)
	defaultPlanningUnitArea   = float64(25)
	defaultRiparianBufferArea = float64(10)
)

func buildExpectedPartialSedimentContribution(params parameters.Parameters) float64 {
	streamDetail := params.GetFloat64(parameters.WaterDensity) * params.GetFloat64(parameters.LocalAcceleration) *
		defaultBankFullFlow * defaultRiverSlope
	adjustedRiparianVegetation := float64(1)
	floodPlainWidthRelationship := 1 - math.Exp(-1.5*math.Pow(10, -2)*defaultFloodPlainWidth)

	maxRiverBankErosion := params.GetFloat64(parameters.BankErosionFudgeFactor) * streamDetail *
		adjustedRiparianVegetation * floodPlainWidthRelationship

	riverSediment := defaultRiverLength * defaultBankHeight * params.GetFloat64(parameters.SedimentDensity) *
		params.GetFloat64(parameters.SuspendedSedimentProportion)

	contribution := maxRiverBankErosion * riverSediment

	return contribution
}

func TestBankSedimentContribution_Initialise(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testDataTable := buildTestTable()
	dummyParameters := new(parameters.Parameters).Initialise()
	contributionUnderTest := new(BankSedimentContribution)

	// when
	actualColumns, actualRows := testDataTable.ColumnAndRowSize()

	// then
	g.Expect(actualRows).To(BeNumerically(equalTo, expectedRowNumber))
	g.Expect(actualColumns).To(BeNumerically(equalTo, expectedColumnNumber))

	// when
	contributionUnderTest.Initialise(testDataTable, *dummyParameters)

	// then
	g.Expect(len(contributionUnderTest.contributionMap)).To(BeNumerically(equalTo, actualRows))

	firstMapEntry := contributionUnderTest.contributionMap["0"]
	expectedPartialSedimentContribution := buildExpectedPartialSedimentContribution(*dummyParameters)
	g.Expect(firstMapEntry.partialSedimentContribution).To(BeNumerically(equalTo, expectedPartialSedimentContribution))
	g.Expect(firstMapEntry.originalIntactRiparianVegetation).To(BeNumerically(equalTo, expectedRiparianVegetationProportion))
}

func TestBankSedimentContribution_OriginalSedimentContribution(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testDataTable := buildTestTable()
	dummyParameters := new(parameters.Parameters).Initialise()
	contributionUnderTest := new(BankSedimentContribution)
	contributionUnderTest.Initialise(testDataTable, *dummyParameters)

	// when
	actualOriginalSedimentContribution := contributionUnderTest.OriginalSedimentContribution()

	// then
	partialSedimentContribution := contributionUnderTest.contributionMap["0"].partialSedimentContribution
	fullSedimentContribution := partialSedimentContribution * contributionUnderTest.adjustedProportionOfIntactVegetation(expectedRiparianVegetationProportion)
	expectedFullSedimentContribution := fullSedimentContribution * expectedRowNumber
	g.Expect(actualOriginalSedimentContribution).To(BeNumerically(equalTo, expectedFullSedimentContribution))
}

func buildTestTable() tables.CsvTable {
	newTable := new(tables.CsvTableImpl)
	newTable.SetColumnAndRowSize(expectedColumnNumber, expectedRowNumber)

	for currentRow := uint(0); currentRow < expectedRowNumber; currentRow++ {
		newTable.SetCell(planningUnitIndex, currentRow, float64(currentRow))
		newTable.SetCell(riverLengthIndex, currentRow, defaultRiverLength)
		newTable.SetCell(channelSlopeIndex, currentRow, defaultRiverSlope)
		newTable.SetCell(bankHeightIndex, currentRow, defaultBankHeight)
		newTable.SetCell(floodPlainWidthIndex, currentRow, defaultFloodPlainWidth)
		newTable.SetCell(bankFullFlowIndex, currentRow, defaultBankFullFlow)
		newTable.SetCell(riparianVegetationIndex, currentRow, expectedRiparianVegetationProportion)
		newTable.SetCell(subCatchmentAreaIndex, currentRow, defaultPlanningUnitArea)
		newTable.SetCell(riparianBufferAreaIndex, currentRow, defaultRiparianBufferArea)
	}

	return newTable
}
