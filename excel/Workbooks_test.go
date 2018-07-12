// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"os"
	"path/filepath"
	"testing"
	. "github.com/onsi/gomega"
)

func TestWorkbooks_Add(t *testing.T) {
	g := NewGomegaWithT(t)

	workbooksUnderTest := excelHandlerUnderTest.Workbooks()

	originalWorkbookCount := workbooksUnderTest.Count()
	g.Expect(originalWorkbookCount).To(BeNumerically("==", 1),"Original Workbooks count should be 1")

	workbook, err := workbooksUnderTest.Add()

	g.Expect(workbook).To(Not(BeNil()),"Workbooks add should return new workbook")
	g.Expect(err).To(BeNil(),"Workbooks Add should not error")

	newWorkbookCount := workbooksUnderTest.Count()

	g.Expect(newWorkbookCount).To(BeNumerically("==", 2),"Workbooks add should increment count")
}

func TestWorkbooks_Open_Bad(t *testing.T) {
	g := NewGomegaWithT(t)

	workbooksUnderTest := excelHandlerUnderTest.Workbooks()
	workbook, err := workbooksUnderTest.Open("badPath")

	g.Expect(workbook).To(BeNil(),"Open Workbooks to bad file path should return null workbook")
	g.Expect(err).To(Not(BeNil()),"Open Workbooks to bad file path should error")
}

func TestWorkbooks_Open_Good(t *testing.T) {
	g := NewGomegaWithT(t)

	workingDirectory, _ := os.Getwd()
	testFixtureAbsolutePath := filepath.Join(workingDirectory, "testdata", "ExcelTestFixture.xls")
	workbooksUnderTest := excelHandlerUnderTest.Workbooks()

	validWorkbook, validWorkbookErr := workbooksUnderTest.Open(testFixtureAbsolutePath)

	g.Expect(validWorkbook).To(Not(BeNil()),"Open Workbooks to good file path should return workbook")
	g.Expect(validWorkbookErr).To(BeNil(),"Open Workbooks to good file path should not error")

	defer validWorkbook.Close()
}
