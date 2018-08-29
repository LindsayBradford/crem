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
	excelHandler     *excel.Handler
	workbook         excel.Workbook
	absoluteFilePath string
	oleWrapper       func(f func())
}

func (eda *ExcelDataAdapter) Initialise() *ExcelDataAdapter {
	defer func() {
		if r := recover(); r != nil {
			eda.excelHandler.Destroy()
			panic(errors.New("Failed initialising via Excel data-source"))
		}
	}()

	eda.excelHandler = new(excel.Handler).Initialise()

	return eda
}

func (eda *ExcelDataAdapter) WithOleFunctionWrapper(wrapper func(f func())) *ExcelDataAdapter {
	eda.oleWrapper = wrapper
	return eda
}

func (eda *ExcelDataAdapter) initialiseDataSource(filePath string) {
	workingDirectory, _ := os.Getwd()
	eda.absoluteFilePath = filepath.Join(workingDirectory, filePath)

	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("could not open data source [" + filePath + "] for initialisation."))
		}
	}()

	eda.oleWrapper(func() {
		workbooks := eda.excelHandler.Workbooks()
		defer workbooks.Release()
		eda.workbook = workbooks.Open(eda.absoluteFilePath)
	})
}

func (eda *ExcelDataAdapter) retrieveAnnealingTableFromWorkbook() (table *annealingTable) {
	table = new(annealingTable)

	eda.oleWrapper(func() {
		worksheet := eda.workbook.WorksheetNamed(data)
		defer worksheet.Release()

		const headerRowCount = uint(1)
		worksheetRowCount := excel.RowCount(worksheet)
		table.rows = make([]annealingData, worksheetRowCount-headerRowCount)

		for index := 0; index < len(table.rows); index++ {
			rowOffset := uint(2 + index)

			costCell := worksheet.Cells(rowOffset, 2)
			table.rows[index].Cost = costCell.Value().(float64)
			costCell.Release()

			featureCell := worksheet.Cells(rowOffset, 3)
			table.rows[index].Feature = featureCell.Value().(float64)
			featureCell.Release()

			puCell := worksheet.Cells(rowOffset, 6)
			table.rows[index].PlanningUnitStatus = (InclusionStatus)(puCell.Value().(float64))
			puCell.Release()
		}
	})

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
	eda.oleWrapper(func() {
		worksheet := eda.workbook.WorksheetNamed(data)
		defer worksheet.Release()
		for index := 0; index < len(table.rows); index++ {
			storageCell := worksheet.Cells(2+uint(index), 5)
			storageCell.SetValue(uint64(table.rows[index].PlanningUnitStatus))
			storageCell.Release()
		}
		excel.AutoFitColumns(worksheet)
	})
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
	eda.oleWrapper(func() {
		worksheet := eda.workbook.WorksheetNamed(tracker)
		defer worksheet.Release()
		excel.ClearUsedRange(worksheet)
	})
}

func (eda *ExcelDataAdapter) storeTrackingTableToWorkbook(table *trackingTable) {
	eda.oleWrapper(func() {

		defer func() {
			if r := recover(); r != nil {
				eda.excelHandler.Destroy()
				panic(errors.New("failed storing data to excel data-source [" + eda.absoluteFilePath + "]"))
			}
		}()

		worksheet := eda.workbook.WorksheetNamed(tracker)
		defer worksheet.Release()
		storeTrackingTableToWorksheet(table, worksheet)
		excel.AutoFitColumns(worksheet)
	})
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

func (eda *ExcelDataAdapter) saveAndCloseWorkbookAs(filePath string) {
	eda.oleWrapper(func() {

		defer func() {
			if r := recover(); r != nil {
				eda.excelHandler.Destroy()
				panic(errors.New("failed saving data to Excel data-source [" + filePath + "]"))
			}
		}()
		eda.workbook.SaveAs(filePath)
		eda.workbook.SetProperty("Saved", true)
		eda.workbook.Close(false)
	})
}

func (eda *ExcelDataAdapter) destroyExcelHandler() {
	eda.oleWrapper(func() {
		eda.excelHandler.Destroy()
	})
}
