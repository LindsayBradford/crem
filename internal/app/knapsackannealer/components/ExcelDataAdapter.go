// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"
	"path/filepath"

	"github.com/LindsayBradford/crm/excel"
)

const tracker = "Tracker"
const data = "Data"

var (
	excelHandler *excel.ExcelHandler
	workbook *excel.Workbook
	testFixtureAbsolutePath string
)

func init() {
	defer func() {
		if r := recover(); r != nil {
			excelHandler.Destroy()
			panic("Failed initialising via Excel data-source")
		}
	}()

	excelHandler = excel.InitialiseHandler()
}

func initialiseDataSource() (filePath string) {
	workingDirectory, _ := os.Getwd()
	testFixtureAbsolutePath = filepath.Join(workingDirectory, "testdata", "KnapsackAnnealerTestFixture.xls")

	var dataSourceErr error
	workbook, dataSourceErr = excelHandler.Workbooks().Open(testFixtureAbsolutePath)

	if dataSourceErr != nil {
		panic("Workbook [" + testFixtureAbsolutePath + "] could not be opened.")
	}

	return testFixtureAbsolutePath
}

func destroyExcelHandler() {
	excelHandler.Destroy()
}

func retrieveAnnealingTableFromWorkbook() (table *annealingTable) {
	table = new(annealingTable)
	worksheet := workbook.WorksheetNamed(data)

	const headerRowCount = uint(1)
	worksheetRowCount := worksheet.UsedRange().Rows().Count()
	table.rows = make([]annealingData, worksheetRowCount - headerRowCount)

	for index := 0; index < len(table.rows); index++ {
		rowOffset := uint(2+index)
		table.rows[index].Cost = worksheet.Cells(rowOffset, 2).Value().(float64)
		table.rows[index].Feature = worksheet.Cells(rowOffset, 3).Value().(float64)
		table.rows[index].PlanningUnitStatus = (InclusionStatus)(worksheet.Cells(rowOffset, 6).Value().(float64))
	}

	randomiseInitialSolutionSet(table)
	return
}

func randomiseInitialSolutionSet(table *annealingTable) {
	for index := 0; index < len(table.rows); index++ {
		randomInOutValue := generateRandomInOutValue()
		table.setPlanningUnitStatusAtIndex(randomInOutValue, uint64(index))
	}
}

func generateRandomInOutValue() InclusionStatus {
	return (InclusionStatus)(randomNumberGenerator.Intn(2))
}

func storeAnnealingTableToWorkbook(table *annealingTable) {
	worksheet := workbook.WorksheetNamed(data)
	for index := 0; index < len(table.rows); index++ {
		worksheet.Cells(2+uint(index), 5).SetValue(uint64(table.rows[index].PlanningUnitStatus))
	}
	worksheet.UsedRange().Columns().AutoFit()
}

func initialiseTrackingTable() *trackingTable {
	clearTrackingDataFromWorkbook()
	return createNewTrackingTable()
}

func createNewTrackingTable() *trackingTable {
	newTrackingTable := new(trackingTable)
	newTrackingTable.headings = []trackingTableHeadings{
		ObjFuncChange,
		Temperature,
		ChangeIsDesirable,
		AcceptanceProbability,
		ChangeAccepted,
		InFirst50,
		InSecond50,
		TotalCost,
	}

	return newTrackingTable
}

func clearTrackingDataFromWorkbook() () {
	worksheet := workbook.WorksheetNamed(tracker)
	worksheet.UsedRange().Clear()
}

func storeTrackingTableToWorkbook(table *trackingTable) {
	defer func() {
		if r := recover(); r != nil {
			excelHandler.Destroy()
			panic("Failed storing data to Excel data-source [" + testFixtureAbsolutePath + "]")
		}
	}()

	worksheet := workbook.WorksheetNamed(tracker)
	storeTrackingTableToWorksheet(table, worksheet)
	worksheet.UsedRange().Columns().AutoFit()
}

func storeTrackingTableToWorksheet(table *trackingTable, worksheet *excel.Worksheet) {
	setTrackingDataColumnHeaders(table, worksheet)
	storeTrackingTableRowsToWorksheet(table, worksheet)
}

func setTrackingDataColumnHeaders(table *trackingTable, worksheet *excel.Worksheet) {
	const headerRowIndex = uint(1)
	for _, heading := range table.headings {
		worksheet.Cells(headerRowIndex, uint(heading)).SetValue(heading.String())
	}
}

func storeTrackingTableRowsToWorksheet(table *trackingTable, worksheet *excel.Worksheet) {
	const rowOffset = 2
	for index := 0; index < len(table.rows); index++ {
		rowNumber := uint(index + rowOffset)
		worksheet.Cells(rowNumber, ObjFuncChange.Index()).SetValue(table.rows[index].ObjectiveFunctionChange)
		worksheet.Cells(rowNumber, Temperature.Index()).SetValue(table.rows[index].Temperature)
		worksheet.Cells(rowNumber, ChangeIsDesirable.Index()).SetValue(table.rows[index].ChangeIsDesirable)
		worksheet.Cells(rowNumber, AcceptanceProbability.Index()).SetValue(table.rows[index].AcceptanceProbability)
		worksheet.Cells(rowNumber, ChangeAccepted.Index()).SetValue(table.rows[index].ChangeAccepted)
		worksheet.Cells(rowNumber, InFirst50.Index()).SetValue(table.rows[index].InFirst50)
		worksheet.Cells(rowNumber, InSecond50.Index()).SetValue(table.rows[index].InSecond50)
		worksheet.Cells(rowNumber, TotalCost.Index()).SetValue(table.rows[index].TotalCost)
	}
}

func saveAndCloseWorkbook() {
	defer func() {
		if r := recover(); r != nil {
			excelHandler.Destroy()
			panic("Failed saving data to Excel data-source [" + testFixtureAbsolutePath + "]")
		}
	}()

	workbook.Save()
	workbook.Close()
}