// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestWorksheets_Add(t *testing.T) {
	g := NewGomegaWithT(t)

	workbooksUnderTest := excelHandlerUnderTest.Workbooks()

	workbook, _ := workbooksUnderTest.Add()
	worksheets, _ := workbook.Worksheets()

	originalWorksheetCount := worksheets.Count()

	g.Expect(originalWorksheetCount).To(BeNumerically("==", uint(1)),"Original Worksheets count should be 1")

	_, worksheetsErr := worksheets.Add()

	g.Expect(worksheetsErr).To(BeNil(),"Worksheets Add should not error")
	newWorksheetCount := worksheets.Count()

	g.Expect(newWorksheetCount).To(BeNumerically("==", originalWorksheetCount + uint(1)),"Worksheets add should increment count")
}