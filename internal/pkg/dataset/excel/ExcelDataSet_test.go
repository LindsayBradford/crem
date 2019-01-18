// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package excel

import (
	"os"
	"path/filepath"
	"testing"

	tables2 "github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	. "github.com/onsi/gomega"
)

func TestDataSet_NewDataSet(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedName := "expectedName"

	dataSetUnderTest := NewDataSet(expectedName, callOnMainThread)

	g.Expect(dataSetUnderTest.Name()).To(BeIdenticalTo(expectedName), "new dataset should have name supplied")
	g.Expect(dataSetUnderTest.Tables()).To(BeEmpty(), "new dataset should have an empty table map")
}

func TestDataSet_Load(t *testing.T) {
	g := NewGomegaWithT(t)

	workingDirectory, _ := os.Getwd()
	testFixtureAbsolutePath := filepath.Join(workingDirectory, "testdata", "testExcelDataSet.xlsx")
	dataSetUnderTest := NewDataSet("testExcelDataSet", callOnMainThread)

	var loadError error
	loadDataSetCall := func() {
		loadError = dataSetUnderTest.Load(testFixtureAbsolutePath)
	}

	g.Expect(loadDataSetCall).To(Not(Panic()), "DataSet Load of good file path should not panic")
	g.Expect(loadError).To(BeNil(), "DataSet Load  to good file path should not return an error ")
	g.Expect(dataSetUnderTest.Tables()).To(Not(BeNil()), "DataSet Load to good file path should return tables")

	tables := dataSetUnderTest.Tables()
	g.Expect(tables).To(HaveKey("testAscTable"), "Loaded dataset has table 'testAscTable'")

	testAscTable := dataSetUnderTest.Tables()["testAscTable"]

	g.Expect(testAscTable.Cell(1, 1)).To(BeNumerically("==", 1), "Loaded dataset has table 'testAscTable'")
	g.Expect(testAscTable.Cell(2, 2)).To(BeNumerically("==", 5), "Loaded dataset has table 'testAscTable'")
	g.Expect(testAscTable.Cell(3, 3)).To(BeNumerically("==", 9), "Loaded dataset has table 'testAscTable'")

	actualAscCols, actualAscRows := testAscTable.Size()
	g.Expect(actualAscCols).To(BeNumerically("==", 5))
	g.Expect(actualAscRows).To(BeNumerically("==", 5))

	g.Expect(tables).To(HaveKey("testCsvTable"), "Loaded dataset has table 'testCsvTable'")

	testCsvTable := dataSetUnderTest.Tables()["testCsvTable"]
	typedCsvTable, _ := testCsvTable.(*tables2.CsvTable)
	g.Expect(typedCsvTable.Header()).To(ContainElement("StringColumn"), "Loaded dataset has table 'testCsvTable'")

	g.Expect(typedCsvTable.Cell(0, 0)).To(BeNumerically("==", 1), "Loaded dataset has table 'testAscTable'")
	g.Expect(typedCsvTable.Cell(1, 1)).To(BeIdenticalTo("entry2"), "Loaded dataset has table 'testAscTable'")
	g.Expect(typedCsvTable.Cell(2, 2)).To(BeNumerically("==", 3.001), "Loaded dataset has table 'testAscTable'")
	g.Expect(typedCsvTable.Cell(3, 3)).To(BeFalse(), "Loaded dataset has table 'testAscTable'")

	actualCsvCols, actualCsvRows := typedCsvTable.Size()
	g.Expect(actualCsvCols).To(BeNumerically("==", 4))
	g.Expect(actualCsvRows).To(BeNumerically("==", 5))
}
