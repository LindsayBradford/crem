// Copyright (c) 2018 Australian Rivers Institute.

package dataset

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	. "github.com/onsi/gomega"
)

func TestDataSet_NewDataSet(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedName := "expectedName"

	dataSetUnderTest := NewDataSet(expectedName)

	g.Expect(dataSetUnderTest.Name()).To(BeIdenticalTo(expectedName), "new dataset should have name supplied")
	g.Expect(dataSetUnderTest.Tables()).To(BeEmpty(), "new dataset should have an empty table map")
}

func TestDataSetImpl_AddTable(t *testing.T) {
	g := NewGomegaWithT(t)

	dataSetUnderTest := NewDataSet("")

	newTable1 := tables.DefaultNullTable
	dataSetUnderTest.AddTable("firstTable", newTable1)

	g.Expect(len(dataSetUnderTest.Tables())).To(BeNumerically("==", 1), "table size should be 1")
	g.Expect(dataSetUnderTest.Table("firstTable")).To(BeIdenticalTo(tables.DefaultNullTable), "added table should be default null table")

	newTable2 := new(tables.CsvTable)
	dataSetUnderTest.AddTable("secondTable", newTable2)

	g.Expect(len(dataSetUnderTest.Tables())).To(BeNumerically("==", 2), "table size should be 2")
	g.Expect(dataSetUnderTest.Table("secondTable")).To(BeIdenticalTo(newTable2), "added table should be new CsvTable")

	newTable3 := new(tables.AscTable)
	dataSetUnderTest.AddTable("thirdTable", newTable3)

	g.Expect(len(dataSetUnderTest.Tables())).To(BeNumerically("==", 3), "table size should be 3")
	g.Expect(dataSetUnderTest.Table("thirdTable")).To(BeIdenticalTo(newTable3), "added table should be new AscTable")
}

func TestDataSetImpl_AddExistingTable_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	dataSetUnderTest := NewDataSet("")

	newTable1 := tables.DefaultNullTable
	firstAddError := dataSetUnderTest.AddTable("firstTable", newTable1)

	g.Expect(firstAddError).To(BeNil(), "addition of non-existent table should succeed")
	g.Expect(len(dataSetUnderTest.Tables())).To(BeNumerically("==", 1), "table size should be 1")
	g.Expect(dataSetUnderTest.Table("firstTable")).To(BeIdenticalTo(tables.DefaultNullTable), "added table should be default null table")

	newTable2 := new(tables.CsvTable)
	secondAddError := dataSetUnderTest.AddTable("firstTable", newTable2)

	g.Expect(secondAddError).To(Not(BeNil()), "addition of pre-existing table should error")
	t.Log(secondAddError)

	g.Expect(len(dataSetUnderTest.Tables())).To(BeNumerically("==", 1), "table size should be 2")
}

func TestDataSetImpl_NonExistentTable_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	dataSetUnderTest := NewDataSet("")

	table, tableError := dataSetUnderTest.Table("notValidTableName")

	g.Expect(table).To(BeNil(), "asking for non-existent table should return nil")
	g.Expect(tableError).To(Not(BeNil()), "asking for non-existent table should raise error")
}
