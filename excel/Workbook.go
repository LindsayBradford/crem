// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
)

type Workbook struct {
	dispatch *ole.IDispatch
}

func (wb *Workbook) Worksheets() (worksheets *Worksheets) {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot retrieve worksheets"))
			worksheets = nil
		}
	}()

	worksheets = new(Worksheets)
	worksheets.dispatch = wb.getProperty("Worksheets")
	return worksheets
}

func (wb *Workbook) getProperty(propertyName string, parameters ...interface{}) *ole.IDispatch {
	return getProperty(wb.dispatch, propertyName, parameters...)
}

func (wb *Workbook) Worksheet(index uint) *Worksheet {
	defer func() {
		if r := recover(); r != nil {
			panic("cannot open worksheet at index [" + fmt.Sprintf("%d", index) + "]")
		}
	}()

	worksheet := new(Worksheet)
	worksheet.dispatch = wb.getProperty("Worksheets", index)
	return worksheet
}

func (wb *Workbook) WorksheetNamed(name string) *Worksheet {
	defer func() {
		if r := recover(); r != nil {
			panic("cannot open worksheet [" + name + "]")
		}
	}()

	worksheet := new(Worksheet)
	worksheet.dispatch = wb.getProperty("Worksheets", name)
	return worksheet
}

func (wb *Workbook) Save() {
	wb.call("Save")
}

func (wb *Workbook) SaveAs(newFileName string) {
	wb.call("SaveAs", newFileName)
}

func (wb *Workbook) Close() {
	wb.call("Close")
}

func (wb *Workbook) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(wb.dispatch, methodName, parameters...)
}
