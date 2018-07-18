// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"
	"path/filepath"

	"github.com/LindsayBradford/crm/excel"
)

var (
	excelHandler *excel.ExcelHandler
	workbook *excel.Workbook
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
	testFixtureAbsolutePath := filepath.Join(workingDirectory, "testdata", "KnapsackAnnealerTestFixture.xls")

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
	worksheet := workbook.WorksheetNamed("Data")

	const headerRowCount = uint(1)
	worksheetRowCount := worksheet.UsedRange().Rows().Count()
	table.rows = make([]annealingData, worksheetRowCount - headerRowCount)

	for index := 0; index < len(table.rows); index++ {
		rowOffset := uint(2+index)
		table.rows[index].Cost = worksheet.Cells(rowOffset, 2).Value().(float64)
		table.rows[index].Feature = worksheet.Cells(rowOffset, 3).Value().(float64)
		table.rows[index].PUStatus = (InclusionStatus)(worksheet.Cells(rowOffset, 6).Value().(float64))
	}

	randomiseInitialSolutionSet(table)
	return
}

func randomiseInitialSolutionSet(table *annealingTable) {
	for index := 0; index < len(table.rows); index++ {
		randomInOutValue := generateRandomInOutValue()
		table.setPUStatusAtIndex(randomInOutValue, uint64(index))
	}
}

func generateRandomInOutValue() InclusionStatus {
	return (InclusionStatus)(randomNumberGenerator.Intn(2))
}

func storeAnnealingTableToWorkbook(table *annealingTable) {
	worksheet := workbook.WorksheetNamed("Data")
	for index := 0; index < len(table.rows); index++ {
		worksheet.Cells(2+uint(index), 5).SetValue(uint64(table.rows[index].PUStatus))
	}
	worksheet.UsedRange().Columns().AutoFit()
}

func clearTrackingDataFromWorkbook() (table *trackingTable) {
	worksheet := workbook.WorksheetNamed("Tracker")
	worksheet.UsedRange().Clear()
	return new(trackingTable)
}

func storeTrackingTableToWorkbook(table *trackingTable) {
	worksheet := workbook.WorksheetNamed("Tracker")
	setTrackingDataColumnHeaders(worksheet)
	storeTrackingTableToWorksheet(table, worksheet)
	worksheet.UsedRange().Columns().AutoFit()
}

func setTrackingDataColumnHeaders(worksheet *excel.Worksheet) {
	columnNames := [...]string{
		"ObjFuncChange",
		"Temperature",
		"ChangeIsDesirable",
		"AcceptanceProbability",
		"ChangeAccepted",
		"InFirst50",
		"InSecond50",
		"TotalCost",
	}

	const headerRowIndex = 1

	for columnIndex := uint(1); columnIndex <= uint(len(columnNames)); columnIndex++ {
		worksheet.Cells(headerRowIndex, columnIndex).SetValue(columnNames[columnIndex-1])
	}
}

func storeTrackingTableToWorksheet(table *trackingTable, worksheet *excel.Worksheet) {
	const rowOffset = 2
	for index := 0; index < len(table.rows); index++ {
		rowNumber := uint(index + rowOffset)
		worksheet.Cells(rowNumber, 1).SetValue(table.rows[index].ObjectiveFunctionChange)
		worksheet.Cells(rowNumber, 2).SetValue(table.rows[index].Temperature)
		worksheet.Cells(rowNumber, 3).SetValue(table.rows[index].ChangeIsDesirable)
		worksheet.Cells(rowNumber, 4).SetValue(table.rows[index].AcceptanceProbability)
		worksheet.Cells(rowNumber, 5).SetValue(table.rows[index].ChangeAccepted)
		worksheet.Cells(rowNumber, 6).SetValue(table.rows[index].InFirst50)
		worksheet.Cells(rowNumber, 7).SetValue(table.rows[index].InSecond50)
		worksheet.Cells(rowNumber, 8).SetValue(table.rows[index].TotalCost)
	}
}

func saveAndCloseWorkbook() {
	workbook.Save()
	workbook.Close()
}