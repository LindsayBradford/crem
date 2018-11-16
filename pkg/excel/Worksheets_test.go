// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestWorksheets_Add(t *testing.T) {
	g := NewGomegaWithT(t)

	workbooksUnderTest := excelHandlerUnderTest.Workbooks()
	defer workbooksUnderTest.Release()

	workbook := workbooksUnderTest.Add()
	defer workbook.Close()
	worksheets := workbook.Worksheets()
	defer worksheets.Release()

	originalWorksheetCount := worksheets.Count()

	g.Expect(originalWorksheetCount).To(BeNumerically("==", uint(1)), "Original Worksheets count should be 1")

	var newWorksheet Worksheet
	addWorksheetCall := func() {
		newWorksheet = worksheets.Add()
	}

	g.Expect(addWorksheetCall).To(Not(Panic()), "Worksheets Add should not panic")
	newWorksheetCount := worksheets.Count()

	g.Expect(newWorksheetCount).To(BeNumerically("==", originalWorksheetCount+uint(1)), "Worksheets add should increment count")
	newWorksheet.Release()
}
