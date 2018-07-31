// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
)

type Worksheets interface {
	Add() (worksheet Worksheet)
	Count() uint
	Item(index uint) Worksheet
}

type WorksheetsImpl struct {
	dispatch *ole.IDispatch
}

func (sheets *WorksheetsImpl) Add() (worksheet Worksheet) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Cannot create new excel worksheet: %s", r)
			panic(errors.New(msg))
			worksheet = nil
		}
	}()

	newWorksheet := new(WorksheetImpl)
	newWorksheet.dispatch = sheets.call("Add")
	return newWorksheet
}

func (sheets *WorksheetsImpl) Count() uint {
	return (uint)(sheets.getPropertyValue("Count"))
}

func (sheets *WorksheetsImpl) Item(index uint) Worksheet {
	worksheetAtIndex := new(WorksheetImpl)
	worksheetAtIndex.dispatch = sheets.getProperty("Item", index)
	return worksheetAtIndex
}

func (sheets *WorksheetsImpl) getProperty(propertyName string, parameters ...interface{}) *ole.IDispatch {
	return getProperty(sheets.dispatch, propertyName, parameters...)
}

func (sheets *WorksheetsImpl) getPropertyValue(propertyName string, parameters ...interface{}) int64 {
	return getPropertyValue(sheets.dispatch, propertyName, parameters...)
}

func (sheets *WorksheetsImpl) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(sheets.dispatch, methodName, parameters...)
}
