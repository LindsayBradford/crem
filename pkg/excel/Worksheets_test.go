// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"testing"

	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestWorksheets_Add(t *testing.T) {
	g := NewGomegaWithT(t)

	// given

	workbooksUnderTest := excelHandlerUnderTest.Workbooks()
	defer workbooksUnderTest.Release()

	workbook := workbooksUnderTest.Add()
	defer workbook.Close()

	worksheets := workbook.Worksheets()
	defer worksheets.Release()

	originalWorksheetCount := worksheets.Count()

	// when

	var newWorksheet Worksheet
	addWorksheetCall := func() {
		newWorksheet = worksheets.Add()
	}

	// then

	expectedWorksheetCount := originalWorksheetCount + 1

	g.Expect(addWorksheetCall).To(Not(Panic()))

	newWorksheetCount := worksheets.Count()

	g.Expect(newWorksheetCount).To(BeNumerically(equalTo, expectedWorksheetCount))
	newWorksheet.Release()
}
