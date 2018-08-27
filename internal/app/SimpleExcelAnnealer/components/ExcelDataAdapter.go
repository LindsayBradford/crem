// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/LindsayBradford/crm/excel"
	"github.com/pkg/errors"
)

const tracker = "Tracker"
const data = "Data"

type ExcelDataAdapter struct {
	excelHandler     *excel.ExcelHandler
	workbook         excel.Workbook
	absoluteFilePath string
}

func (eda *ExcelDataAdapter) Initialise() {
	defer func() {
		if r := recover(); r != nil {
			eda.excelHandler.Destroy()
			panic(errors.New("Failed initialising via Excel data-source"))
		}
	}()

	eda.excelHandler = excel.InitialiseHandler()
}

func (eda *ExcelDataAdapter) initialiseDataSource(filePath string) {
	workingDirectory, _ := os.Getwd()
	eda.absoluteFilePath = filepath.Join(workingDirectory, filePath)

	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("could not open data source [" + filePath + "] for initialisation."))
		}
	}()

	eda.workbook = eda.excelHandler.Workbooks().Open(eda.absoluteFilePath)
}

func (eda *ExcelDataAdapter) destroyExcelHandler() {
	eda.excelHandler.Destroy()
}

func (eda *ExcelDataAdapter) retrieveAnnealingTableFromWorkbook() (table *annealingTable) {
	table = new(annealingTable)
	worksheet := eda.workbook.WorksheetNamed(data)

	const headerRowCount = uint(1)
	worksheetRowCount := worksheet.UsedRange().Rows().Count()
	table.rows = make([]annealingData, worksheetRowCount-headerRowCount)

	for index := 0; index < len(table.rows); index++ {
		rowOffset := uint(2 + index)
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

func (eda *ExcelDataAdapter) storeAnnealingTableToWorkbook(table *annealingTable) {
	worksheet := eda.workbook.WorksheetNamed(data)
	for index := 0; index < len(table.rows); index++ {
		worksheet.Cells(2+uint(index), 5).SetValue(uint64(table.rows[index].PlanningUnitStatus))
	}
	worksheet.UsedRange().Columns().AutoFit()
}

func (eda *ExcelDataAdapter) initialiseTrackingTable() *trackingTable {
	eda.clearTrackingDataFromWorkbook()
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

func (eda *ExcelDataAdapter) clearTrackingDataFromWorkbook() {
	worksheet := eda.workbook.WorksheetNamed(tracker)
	worksheet.UsedRange().Clear()
}

func (eda *ExcelDataAdapter) storeTrackingTableToWorkbook(table *trackingTable) {
	defer func() {
		if r := recover(); r != nil {
			eda.excelHandler.Destroy()
			panic(errors.New("failed storing data to excel data-source [" + eda.absoluteFilePath + "]"))
		}
	}()

	worksheet := eda.workbook.WorksheetNamed(tracker)
	storeTrackingTableToWorksheet(table, worksheet)
	worksheet.UsedRange().Columns().AutoFit()
}

func storeTrackingTableToWorksheet(table *trackingTable, worksheet excel.Worksheet) {
	tempFile := writeTableToTempCsvFile(table)
	loadTempCsvFileIntoWorksheet(tempFile, worksheet)
	deleteTempCsvFile(tempFile)
}

func writeTableToTempCsvFile(table *trackingTable) string {
	tmpFileHandle, tmpFileErr := ioutil.TempFile("", "trackingTableCsv")
	if tmpFileErr != nil {
		panic(errors.Wrap(tmpFileErr, "failed creating temporary table csv file"))
	}

	w := csv.NewWriter(tmpFileHandle)
	w.UseCRLF = true

	headingsAsStringArray := make([]string, len(table.headings))
	for _, heading := range table.headings {
		headingsAsStringArray[heading.Index()-1] = heading.String()
	}
	if headingWriteError := w.Write(headingsAsStringArray); headingWriteError != nil {
		panic(errors.Wrap(headingWriteError, "failed writing header row to csv file"))
	}

	for rowIndex, row := range table.rows {
		rowArray := make([]string, len(table.headings))

		objectivesFunctionValueAsString := fmt.Sprintf("%f", row.ObjectiveFunctionChange)
		rowArray[ObjFuncChange.Index()-1] = objectivesFunctionValueAsString

		TemperatureAsString := fmt.Sprintf("%f", row.Temperature)
		rowArray[Temperature.Index()-1] = TemperatureAsString

		ChangeIsDesirableAsString := fmt.Sprintf("%t", row.ChangeIsDesirable)
		rowArray[ChangeIsDesirable.Index()-1] = ChangeIsDesirableAsString

		AcceptanceProbabilityAsString := fmt.Sprintf("%f", row.AcceptanceProbability)
		rowArray[AcceptanceProbability.Index()-1] = AcceptanceProbabilityAsString

		ChangeAcceptedAsString := fmt.Sprintf("%t", row.ChangeAccepted)
		rowArray[ChangeAccepted.Index()-1] = ChangeAcceptedAsString

		InFirst50AsString := fmt.Sprintf("%d", row.InFirst50)
		rowArray[InFirst50.Index()-1] = InFirst50AsString

		InSecond50AsString := fmt.Sprintf("%d", row.InSecond50)
		rowArray[InSecond50.Index()-1] = InSecond50AsString

		TotalCostAsString := fmt.Sprintf("%f", row.TotalCost)
		rowArray[TotalCost.Index()-1] = TotalCostAsString

		if writeError := w.Write(rowArray); writeError != nil {
			panic(errors.Wrapf(writeError, "failed writing record [%d} to csv file", rowIndex))
		}
	}

	w.Flush()

	if closeErr := tmpFileHandle.Close(); closeErr != nil {
		panic(errors.Wrap(closeErr, "failed to close csv file"))
	}

	return tmpFileHandle.Name()
}

func loadTempCsvFileIntoWorksheet(tempFileName string, worksheet excel.Worksheet) {
	excel.AddCsvFileContentToWorksheet(tempFileName, worksheet)
}

func deleteTempCsvFile(tempFileName string) {
	defer os.Remove(tempFileName)
}

func (eda *ExcelDataAdapter) saveWorkbookAs(filePath string) {
	eda.workbook.SaveAs(filePath)
}

func (eda *ExcelDataAdapter) saveAndCloseWorkbookAs(filePath string) {
	defer func() {
		if r := recover(); r != nil {
			eda.excelHandler.Destroy()
			panic(errors.New("failed saving data to Excel data-source [" + filePath + "]"))
		}
	}()

	eda.saveWorkbookAs(filePath)
	eda.workbook.Close()
}
