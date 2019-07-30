// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"fmt"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"
)

type Workbook interface {
	Worksheets() (worksheets Worksheets)
	Worksheet(index uint) Worksheet
	WorksheetNamed(name string) Worksheet
	Save()
	SaveAs(newFileName string)
	Close(args ...interface{})
	SetProperty(propertyName string, propertyValue interface{})
	Release()
}

type WorkbookImpl struct {
	oleWrapper
}

func (wb *WorkbookImpl) WithDispatch(dispatch *ole.IDispatch) *WorkbookImpl {
	wb.dispatch = dispatch
	return wb
}

func (wb *WorkbookImpl) Worksheets() (worksheets Worksheets) {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot retrieve worksheets"))
		}
	}()

	dispatch := wb.getProperty("Worksheets")
	return new(WorksheetsImpl).WithDispatch(dispatch)
}

func (wb *WorkbookImpl) Worksheet(index uint) Worksheet {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot open worksheet at index [" + fmt.Sprintf("%d", index) + "]"))
		}
	}()

	dispatch := wb.getProperty("Worksheets", index)
	return new(WorksheetImpl).WithDispatch(dispatch)
}

func (wb *WorkbookImpl) WorksheetNamed(name string) Worksheet {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		panic(errors.New("cannot open worksheet [" + name + "]"))
	// 	}
	// }()
	defer handleGoOleLibraryPanicAsErrorPanic()
	dispatch := wb.getProperty("Worksheets", name)
	return new(WorksheetImpl).WithDispatch(dispatch)
}

func (wb *WorkbookImpl) Save() {
	defer handleGoOleLibraryPanicAsErrorPanic()
	wb.call("Save")
}

func (wb *WorkbookImpl) SaveAs(newFileName string) {
	defer handleGoOleLibraryPanicAsErrorPanic()
	if !filepath.IsAbs(newFileName) {
		panic(errors.New("cannot save-as workbook to relative path [" + newFileName + "]"))
	}

	wb.call("SaveAs", newFileName)
}

func handleGoOleLibraryPanicAsErrorPanic() {
	// No idea why why go-ole are stripping the error wrapper from their panic calls, but they are.
	// Here, we wrap them in an error again so higher-level panic handlers have a standard interface to deal with.
	if r := recover(); r != nil {
		if recoveredString, isString := r.(string); isString {
			recoveredError := errors.New(recoveredString)
			panic(errors.Wrap(recoveredError, "excel workbook"))
		}
		panic(r)
	}
}

func (wb *WorkbookImpl) Close(parameters ...interface{}) {
	wb.call("Close", parameters...)
	wb.Release()
}

func (wb *WorkbookImpl) getProperty(propertyName string, parameters ...interface{}) *ole.IDispatch {
	return getProperty(wb.dispatch, propertyName, parameters...)
}

func (wb *WorkbookImpl) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(wb.dispatch, methodName, parameters...)
}

func (wb *WorkbookImpl) SetProperty(propertyName string, propertyValue interface{}) {
	setProperty(wb.dispatch, propertyName, propertyValue)
}
