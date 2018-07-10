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

	workingDirectory, _ := os.Getwd()
	testFixtureAbsolutePath := filepath.Join(workingDirectory, "testdata", "ExcelTestFixture.xls")
	workbooksUnderTest := excelHandlerUnderTest.Workbooks()

	workbookUnderTest, _ := workbooksUnderTest.Open(testFixtureAbsolutePath)
	worksheetsUnderTest, _ := workbookUnderTest.Worksheets()

	worksheetCount := worksheetsUnderTest.Count()
	g.Expect(worksheetCount).To(BeIdenticalTo(uint(2)), "Expected worksheets count of 2 for test fixture")

	worksheetOne, worksheetOneErr := workbookUnderTest.Worksheet(1)
	g.Expect(worksheetOneErr).To(BeNil(), "Worksheet(1) should not error")
	g.Expect(worksheetOne).To(Not(BeNil()), "Worksheet(1) should not be nil")
	g.Expect(worksheetOne.Name()).To(Equal("FirstSheet"), "Expected worksheet(1) name was 'FirstSheet'")

	worksheetTwo, worksheetTwoErr := workbookUnderTest.Worksheet(2)
	g.Expect(worksheetTwoErr).To(BeNil(), "Worksheet(2) should not error")
	g.Expect(worksheetTwo).To(Not(BeNil()), "Worksheet(2) should not be nil")
	g.Expect(worksheetTwo.Name()).To(Equal("SecondSheet"), "Expected worksheet(2) name was 'SecondSheet'")

	worksheetThree, worksheetThreeErr := workbookUnderTest.Worksheet(3)
	g.Expect(worksheetThreeErr).To(Not(BeNil()), "Worksheet(3) should error")
	g.Expect(worksheetThree).To(BeNil(), "Worksheet(3) should be nil")
}
