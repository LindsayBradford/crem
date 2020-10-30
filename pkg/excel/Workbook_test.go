// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"testing"
)

func TestWorkbook_Worksheet(t *testing.T) {
	g := NewGomegaWithT(t)

	workbooksUnderTest := excelHandlerUnderTest.Workbooks()
	defer workbooksUnderTest.Release()

	workingDirectory, _ := os.Getwd()
	testFixtureAbsolutePath := filepath.Join(workingDirectory, "testdata", "ExcelTestFixture.xls")
	workbookUnderTest := workbooksUnderTest.Open(testFixtureAbsolutePath)
	defer workbookUnderTest.Release()

	worksheetsUnderTest := workbookUnderTest.Worksheets()
	defer worksheetsUnderTest.Release()

	worksheetCount := worksheetsUnderTest.Count()
	g.Expect(worksheetCount).To(BeIdenticalTo(uint(2)), "Expected worksheets count of 2 for test fixture")

	worksheetOne := workbookUnderTest.Worksheet(1)
	defer worksheetOne.Release()
	g.Expect(worksheetOne).To(Not(BeNil()), "Worksheet(1) should not be nil")
	g.Expect(worksheetOne.Name()).To(Equal("FirstSheet"), "Expected worksheet(1) name was 'FirstSheet'")
}
