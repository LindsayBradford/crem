// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/LindsayBradford/crm/excel"
	"github.com/onsi/gomega/gstruct/errors"
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
	testFixtureAbsolutePath = filepath.Join(workingDirectory, "testdata", "SimpleExcelAnnealerTestFixture.xls")

	defer func() {
		if r := recover(); r != nil {
			panic("Workbook [" + testFixtureAbsolutePath + "] could not be opened. Data source initialisation failed.")
		}
	}()

	workbook = excelHandler.Workbooks().Open(testFixtureAbsolutePath)

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
	tempFile := writeTableToTempCsvFile(table)
	loadTempCsvFileIntoWorksheet(tempFile, worksheet)
	deleteTempCsvFile(tempFile)
}

func writeTableToTempCsvFile(table *trackingTable) string {
	tmpFileHandle, err := ioutil.TempFile("", "trackingTableCsv")
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(tmpFileHandle)
	w.UseCRLF = true

	headingsAsStringArray := make([]string, len(table.headings))
	for _, heading := range table.headings {
		headingsAsStringArray[heading.Index() - 1] = heading.String()
	}
	if err := w.Write(headingsAsStringArray); err != nil {
		panic(errors.Nest("error writing header to csv:", err))
	}

  for _, row := range table.rows {
  	rowArray := make([]string, len(table.headings))

  	objectivesFunctionValueAsString := fmt.Sprintf("%f",row.ObjectiveFunctionChange)
	  rowArray[ObjFuncChange.Index() - 1] = objectivesFunctionValueAsString

	  TemperatureAsString := fmt.Sprintf("%f",row.Temperature)
	  rowArray[Temperature.Index() - 1] = TemperatureAsString

	  ChangeIsDesirableAsString := fmt.Sprintf("%t",row.ChangeIsDesirable)
	  rowArray[ChangeIsDesirable.Index() - 1] = ChangeIsDesirableAsString

	  AcceptanceProbabilityAsString := fmt.Sprintf("%f",row.AcceptanceProbability)
	  rowArray[AcceptanceProbability.Index() - 1] = AcceptanceProbabilityAsString

	  ChangeAcceptedAsString := fmt.Sprintf("%t",row.ChangeAccepted)
	  rowArray[ChangeAccepted.Index() - 1] = ChangeAcceptedAsString

	  InFirst50AsString := fmt.Sprintf("%d",row.InFirst50)
	  rowArray[InFirst50.Index() - 1] = InFirst50AsString

	  InSecond50AsString := fmt.Sprintf("%d",row.InSecond50)
	  rowArray[InSecond50.Index() - 1] = InSecond50AsString

	  TotalCostAsString := fmt.Sprintf("%f",row.TotalCost)
	  rowArray[TotalCost.Index() - 1] = TotalCostAsString

		if err := w.Write(rowArray); err != nil {
			panic(errors.Nest("error writing record to csv:", err))
		}
	}

	w.Flush()

	if err := tmpFileHandle.Close(); err != nil {
		panic(errors.Nest("error closing csv:", err))
	}

	return tmpFileHandle.Name()
}

func loadTempCsvFileIntoWorksheet(tempFileName string, worksheet *excel.Worksheet) {
	excel.AddCsvFileContentToWorksheet(tempFileName, worksheet)
}

func deleteTempCsvFile(tempFileName string) {
	defer os.Remove(tempFileName)
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