// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"testing"
)

func TestWorksheet_Cells(t *testing.T) {
	g := NewGomegaWithT(t)

	workingDirectory, _ := os.Getwd()
	testFixtureAbsolutePath := filepath.Join(workingDirectory, "testdata", "ExcelTestFixture.xls")
	workbooksUnderTest := excelHandlerUnderTest.Workbooks()

	workbookUnderTest := workbooksUnderTest.Open(testFixtureAbsolutePath)

	worksheetOne := workbookUnderTest.Worksheet(1)
	worksheetOneActualColumns := worksheetOne.UsedRange().Columns().Count()
	g.Expect(worksheetOneActualColumns).To(BeIdenticalTo(uint(3)), "Columns used should be 3")

	worksheetOneActualRows := worksheetOne.UsedRange().Rows().Count()
	g.Expect(worksheetOneActualRows).To(BeIdenticalTo(uint(4)), "Rows used should be 4")

	cell0_0 := worksheetOne.Cells(0, 0)
	g.Expect(cell0_0).To(BeNil(), "Cells(0,0) should be nil")

	cell1_1 := worksheetOne.Cells(1, 1)
	g.Expect(cell1_1.Value()).To(Equal("PlanningUnit"), "Cells(1,1) value should be 'PlanningUnit'")

	cell1_2 := worksheetOne.Cells(1, 2)
	g.Expect(cell1_2.Value()).To(Equal("Cost"), "Cells(1,2) value should be 'Cost'")

	cell1_3 := worksheetOne.Cells(1, 3)
	g.Expect(cell1_3.Value()).To(Equal("Feature"), "Cells(1,3) value should be 'Feature'")

	cell2_1 := worksheetOne.Cells(2, 1)
	g.Expect(cell2_1.Value()).To(BeNumerically("==", 1), "Cells(2,1) value should be 1")

	cell2_2 := worksheetOne.Cells(2, 2)
	g.Expect(cell2_2.Value()).To(BeNumerically("~", 2.1163097067102, 1e-13), "Cells(2,2) value should ~= 2.1163097067102")

	cell2_3 := worksheetOne.Cells(2, 3)
	g.Expect(cell2_3.Value()).To(BeNumerically("~", 2.11999215931333, 1e-13), "Cells(2,3) value should ~= 2.11999215931333")

	expectedCell5_5Value := "testCell5_5_value"
	worksheetOne.Cells(5, 5).SetValue(expectedCell5_5Value)
	actualCell5_5Value := worksheetOne.Cells(5, 5).Value()

	g.Expect(actualCell5_5Value).To(Equal(expectedCell5_5Value), "Cells(5,5) actual value match expected value")
}
