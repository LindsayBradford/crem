// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
)

type Workbook interface {
	Worksheets() (worksheets Worksheets)
	Worksheet(index uint) Worksheet
	WorksheetNamed(name string) Worksheet
	Save()
	SaveAs(newFileName string)
	Close()
}

type WorkbookImpl struct {
	dispatch *ole.IDispatch
}

func (wb *WorkbookImpl) Worksheets() (worksheets Worksheets) {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot retrieve worksheets"))
			worksheets = nil
		}
	}()

	newWorksheets := new(WorksheetsImpl)
	newWorksheets.dispatch = wb.getProperty("Worksheets")
	return newWorksheets
}

func (wb *WorkbookImpl) Worksheet(index uint) Worksheet {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot open worksheet at index [" + fmt.Sprintf("%d", index) + "]"))
		}
	}()

	newWorksheet := new(WorksheetImpl)
	newWorksheet.dispatch = wb.getProperty("Worksheets", index)
	return newWorksheet
}

func (wb *WorkbookImpl) WorksheetNamed(name string) Worksheet {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot open worksheet [" + name + "]"))
		}
	}()

	namedWorksheet := new(WorksheetImpl)
	namedWorksheet.dispatch = wb.getProperty("Worksheets", name)
	return namedWorksheet
}

func (wb *WorkbookImpl) Save() {
	wb.call("Save")
}

func (wb *WorkbookImpl) SaveAs(newFileName string) {
	wb.call("SaveAs", newFileName)
}

func (wb *WorkbookImpl) Close() {
	wb.call("Close")
}

func (wb *WorkbookImpl) getProperty(propertyName string, parameters ...interface{}) *ole.IDispatch {
	return getProperty(wb.dispatch, propertyName, parameters...)
}

func (wb *WorkbookImpl) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(wb.dispatch, methodName, parameters...)
}
