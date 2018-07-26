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

	workbook := workbooksUnderTest.Add()
	worksheets := workbook.Worksheets()

	originalWorksheetCount := worksheets.Count()

	g.Expect(originalWorksheetCount).To(BeNumerically("==", uint(1)), "Original Worksheets count should be 1")

	addWorksheetCall := func() {
		worksheets.Add()
	}

	g.Expect(addWorksheetCall).To(Not(Panic()), "Worksheets Add should not panic")
	newWorksheetCount := worksheets.Count()

	g.Expect(newWorksheetCount).To(BeNumerically("==", originalWorksheetCount+uint(1)), "Worksheets add should increment count")
}
