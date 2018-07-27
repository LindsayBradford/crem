// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"testing"
)

func TestUtils_AddWorksheetFromCsvFileToWorkbook(t *testing.T) {
	g := NewGomegaWithT(t)

	workbooksUnderTest := excelHandlerUnderTest.Workbooks()

	workingDirectory, _ := os.Getwd()
	testFixtureAbsolutePath := filepath.Join(workingDirectory, "testdata", "ExcelTestFixture.xls")

	workbook := workbooksUnderTest.Open(testFixtureAbsolutePath)
	worksheets := workbook.Worksheets()

	originalWorksheetCount := worksheets.Count()

	g.Expect(originalWorksheetCount).To(BeNumerically("==", uint(2)), "Original Worksheets count should be 2")

	const csvTestWorksheetName = "csvTestWorksheet"
	addWorksheetCall := func() {
		workingDir, _ := os.Getwd()
		csvTestFixtureAbsolutePath := filepath.Join(workingDir, "testdata", "CSVTestFixture.csv")
		AddWorksheetFromCsvFileToWorkbook(csvTestFixtureAbsolutePath, csvTestWorksheetName, workbook)
	}

	g.Expect(addWorksheetCall).To(Not(Panic()), "Worksheets Add should not panic")
	newWorksheetCount := worksheets.Count()

	g.Expect(newWorksheetCount).To(BeNumerically("==", originalWorksheetCount+uint(1)), "Worksheets add should increment count")

	tempOutputAbsolutePath := filepath.Join(workingDirectory, "testdata", "tempOutput.xls")
	workbook.SaveAs(tempOutputAbsolutePath)

	newWorksheet := workbook.WorksheetNamed(csvTestWorksheetName)

	cell1_1 := newWorksheet.Cells(1, 1)
	g.Expect(cell1_1.Value()).To(Equal("PlanningUnit"), "Cells(1,1) value should be 'PlanningUnit'")

	cell1_2 := newWorksheet.Cells(1, 2)
	g.Expect(cell1_2.Value()).To(Equal("Cost"), "Cells(1,2) value should be 'Cost'")

	cell1_3 := newWorksheet.Cells(1, 3)
	g.Expect(cell1_3.Value()).To(Equal("Feature"), "Cells(1,3) value should be 'Feature'")

	cell2_1 := newWorksheet.Cells(2, 1)
	g.Expect(cell2_1.Value()).To(BeNumerically("==", 7), "Cells(2,1) value should be 7")

	cell2_2 := newWorksheet.Cells(2, 2)
	g.Expect(cell2_2.Value()).To(BeNumerically("~", 2.5, 1e-13), "Cells(2,2) value should ~= 2.5")

	cell2_3 := newWorksheet.Cells(2, 3)
	g.Expect(cell2_3.Value()).To(BeNumerically("~", 2.45, 1e-13), "Cells(2,3) value should ~= 2.45")

	workbook.Close()
	os.Remove(tempOutputAbsolutePath)
}
