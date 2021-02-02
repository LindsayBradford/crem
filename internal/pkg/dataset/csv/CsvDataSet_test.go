// Copyright (c) 2019 Australian Rivers Institute.

package csv

import (
	"io/ioutil"
	"testing"

	tables2 "github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	. "github.com/onsi/gomega"
)

func TestDataSet_NewDataSet(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedName := "expectedName"

	dataSetUnderTest := NewDataSet(expectedName)

	g.Expect(dataSetUnderTest.Name()).To(BeIdenticalTo(expectedName), "new dataset should have name supplied")
	g.Expect(dataSetUnderTest.Tables()).To(BeEmpty(), "new dataset should have an empty table map")
}

func TestDataSet_Load_MissingFile_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testFixturePath := "testdata/missingCsvFile.csv"
	dataSetUnderTest := NewDataSet("testDataSet")

	// when
	var loadError error
	loadError = dataSetUnderTest.Load(testFixturePath)

	// then
	g.Expect(loadError).To(Not(BeNil()), "DataSet Load to bad file path should return error ")
	t.Log(loadError)

	g.Expect(len(dataSetUnderTest.Tables())).To(BeNumerically("==", 0), "DataSet Load to bad file path should return zero tables")
}

func TestDataSet_Load_MalformedMetaFile_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testFixturePath := "testdata/malformedHeadingsFile.csv"
	dataSetUnderTest := NewDataSet("testDataSet")

	// when
	var loadError error
	loadError = dataSetUnderTest.Load(testFixturePath)

	// then
	g.Expect(loadError).To(Not(BeNil()), "DataSet Load of bad headings file should return error ")
	t.Log(loadError)

	g.Expect(len(dataSetUnderTest.Tables())).To(BeNumerically("==", 0), "DataSet Load to bad file path should return zero tables")
}

func TestDataSet_Load_Directory_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testFixturePath := "testdata"
	dataSetUnderTest := NewDataSet("testDataSet")

	// when
	var loadError error
	loadError = dataSetUnderTest.Load(testFixturePath)

	// then
	g.Expect(loadError).To(Not(BeNil()), "DataSet Load of bad headings file should return error ")
	t.Log(loadError)

	g.Expect(len(dataSetUnderTest.Tables())).To(BeNumerically("==", 0), "DataSet Load to bad file path should return zero tables")
}

func TestDataSet_Load_MissingTableFile_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	testFixturePath := "testdata/missingTableFiles.csv"
	dataSetUnderTest := NewDataSet("testDataSet")

	// when
	var loadError error
	loadError = dataSetUnderTest.Load(testFixturePath)

	// then
	g.Expect(loadError).To(Not(BeNil()), "DataSet Load to bad file path should return error ")
	t.Log(loadError)

	g.Expect(len(dataSetUnderTest.Tables())).To(BeNumerically("==", 0), "DataSet Load to bad file path should return zero tables")
}

func TestDataSet_Load_ValidDataSet(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	validDataSetPath := "testdata/validDataSet.csv"
	dataSetUnderTest := NewDataSet("dataSetUnderTest")

	// when
	var loadError error
	loadDataSetCall := func() {
		loadError = dataSetUnderTest.Load(validDataSetPath)
	}

	// then
	g.Expect(loadDataSetCall).To(Not(Panic()), "DataSet Load of good file path should not panic")
	g.Expect(loadError).To(BeNil(), "DataSet Load  to good file path should not return an error ")
	g.Expect(dataSetUnderTest.Tables()).To(Not(BeNil()), "DataSet Load to good file path should return tables")

	tables := dataSetUnderTest.Tables()

	g.Expect(tables).To(HaveKey("CsvTable"))

	testCsvTable := dataSetUnderTest.Tables()["CsvTable"]
	typedCsvTable, _ := testCsvTable.(tables2.CsvTable)
	g.Expect(typedCsvTable.Header()).To(ContainElement("StringColumn"))

	g.Expect(typedCsvTable.Cell(0, 0)).To(BeNumerically("==", 1))
	g.Expect(typedCsvTable.Cell(1, 1)).To(BeIdenticalTo("entry2"))
	g.Expect(typedCsvTable.Cell(2, 2)).To(BeNumerically("==", 3.001))
	g.Expect(typedCsvTable.Cell(3, 3)).To(BeFalse())

	actualCsvCols, actualCsvRows := typedCsvTable.ColumnAndRowSize()
	g.Expect(actualCsvCols).To(BeNumerically("==", 4))
	g.Expect(actualCsvRows).To(BeNumerically("==", 5))
}

func TestDataTable_Parse_Valid(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	dataSetUnderTest := NewDataSet("dataSetUnderTest")
	csvText := loadTextFromFile(g, "testdata/validCsvFile.csv")

	// when
	parseDataSetCall := func() {
		dataSetUnderTest.ParseCsvTextIntoTable("testTable", csvText)
	}

	// then
	g.Expect(parseDataSetCall).To(Not(Panic()), "DataSet Load of good file path should not panic")
	g.Expect(dataSetUnderTest.Errors()).To(BeNil(), "DataSet Load  to good file path should not return an error ")
	g.Expect(dataSetUnderTest.Tables()).To(Not(BeNil()), "DataSet Load to good file path should return tables")

	tables := dataSetUnderTest.Tables()

	g.Expect(tables).To(HaveKey("testTable"))

	testCsvTable := dataSetUnderTest.Tables()["testTable"]
	typedCsvTable, _ := testCsvTable.(tables2.CsvTable)
	g.Expect(typedCsvTable.Header()).To(ContainElement("StringColumn"))

	g.Expect(typedCsvTable.Cell(0, 0)).To(BeNumerically("==", 1))
	g.Expect(typedCsvTable.Cell(1, 1)).To(BeIdenticalTo("entry2"))
	g.Expect(typedCsvTable.Cell(2, 2)).To(BeNumerically("==", 3.001))
	g.Expect(typedCsvTable.Cell(3, 3)).To(BeFalse())

	actualCsvCols, actualCsvRows := typedCsvTable.ColumnAndRowSize()
	g.Expect(actualCsvCols).To(BeNumerically("==", 4))
	g.Expect(actualCsvRows).To(BeNumerically("==", 5))
}

func loadTextFromFile(g *GomegaWithT, filePath string) string {
	fileContent, openError := ioutil.ReadFile(filePath)

	g.Expect(openError).To(BeNil())

	text := string(fileContent)
	return text
}
