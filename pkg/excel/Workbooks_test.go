// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"testing"
)

func TestWorkbooks_Add(t *testing.T) {
	g := NewGomegaWithT(t)

	workbooksUnderTest := excelHandlerUnderTest.Workbooks()
	defer workbooksUnderTest.Release()

	// Warning:  THis is driven by generic user-overridable Excel configuration.
	// This number must match what your locally installed Excel configuration has been set to.
	const expectedInitialValue = 1

	originalWorkbookCount := workbooksUnderTest.Count()
	g.Expect(originalWorkbookCount).To(BeNumerically("==", expectedInitialValue), "Original Workbooks count should be 1")

	var workbook Workbook
	addWWorkbookCall := func() {
		workbook = workbooksUnderTest.Add()
	}

	g.Expect(addWWorkbookCall).To(Not(Panic()), "Workbooks Add should not panic")
	g.Expect(workbook).To(Not(BeNil()), "Workbooks add should return new workbook")

	newWorkbookCount := workbooksUnderTest.Count()

	g.Expect(newWorkbookCount).To(BeNumerically("==", expectedInitialValue+1), "Workbooks add should increment count")
	workbook.Close()
	workbook.Release()
}

func TestWorkbooks_Open_Bad(t *testing.T) {
	g := NewGomegaWithT(t)

	workbooksUnderTest := excelHandlerUnderTest.Workbooks()
	defer workbooksUnderTest.Release()

	var workbook *Workbook
	addWWorkbookCall := func() {
		workbooksUnderTest.Open("badPath")
	}
	g.Expect(addWWorkbookCall).To(Panic(), "Workbooks Add of bad file path should panic")
	g.Expect(workbook).To(BeNil(), "Open Workbooks to bad file path should return null workbook")
}

func TestWorkbooks_Open_Good(t *testing.T) {
	g := NewGomegaWithT(t)

	workingDirectory, _ := os.Getwd()
	testFixtureAbsolutePath := filepath.Join(workingDirectory, "testdata", "ExcelTestFixture.xls")
	workbooksUnderTest := excelHandlerUnderTest.Workbooks()
	defer workbooksUnderTest.Release()

	var validWorkbook Workbook
	addWWorkbookCall := func() {
		validWorkbook = workbooksUnderTest.Open(testFixtureAbsolutePath)
	}

	g.Expect(addWWorkbookCall).To(Not(Panic()), "Workbooks Open of good file path should not panic")
	g.Expect(validWorkbook).To(Not(BeNil()), "Open Workbooks to good file path should return workbook")

	validWorkbook.Close()
}
