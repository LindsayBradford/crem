// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package excel

import (
	"testing"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"

	. "github.com/onsi/gomega"
)

func TestGoOle_NoDanglingProcessIsolationGoOle(t *testing.T) {
	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)

	excelAppObject, _ := oleutil.CreateObject("Excel.Application")
	defer excelAppObject.Release()

	newExcelIDispatch, _ := excelAppObject.QueryInterface(ole.IID_IDispatch)
	newExcelIDispatch.Release()

	ole.CoUninitialize()
	t.Log("Check for running excel processes to ensure there are no dangling processes")
}

func TestExcelWrapper_JustSetupAndTearDown(t *testing.T) {
	t.Log("Check for running excel processes to ensure there are no dangling processes")
}

func TestGoOle_DanglingProcessIsolation(t *testing.T) {
	g := NewGomegaWithT(t)

	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	excelAppObject, _ := oleutil.CreateObject("Excel.Application")

	newExcelIDispatch, _ := excelAppObject.QueryInterface(ole.IID_IDispatch)

	workbooks := oleutil.MustGetProperty(newExcelIDispatch, "Workbooks").ToIDispatch()

	workbookCount := oleutil.MustGetProperty(workbooks, "Count").Value()
	g.Expect(workbookCount).To(BeNumerically("==", 0))

	newWorkbook := oleutil.MustCallMethod(workbooks, "Add", nil).ToIDispatch()

	workbookCount = oleutil.MustGetProperty(workbooks, "Count").Value()
	g.Expect(workbookCount).To(BeNumerically("==", 1))

	oleutil.MustCallMethod(newWorkbook, "Close", false)
	workbooks.Release()

	newExcelIDispatch.Release()
	ole.CoUninitialize()
}
