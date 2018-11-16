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
	Release()
}

type WorksheetsImpl struct {
	dispatch *ole.IDispatch
}

func (sheets *WorksheetsImpl) WithDispatch(dispatch *ole.IDispatch) *WorksheetsImpl {
	sheets.dispatch = dispatch
	return sheets
}

func (sheets *WorksheetsImpl) Add() (worksheet Worksheet) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Cannot create new excel worksheet: %s", r)
			panic(errors.New(msg))
		}
	}()

	dispatch := sheets.call("Add")
	return new(WorksheetImpl).WithDispatch(dispatch)
}

func (sheets *WorksheetsImpl) Count() uint {
	return (uint)(sheets.getPropertyValue("Count"))
}

func (sheets *WorksheetsImpl) Item(index uint) Worksheet {
	dispatch := sheets.getProperty("Item", index)
	return new(WorksheetImpl).WithDispatch(dispatch)
}

func (sheets *WorksheetsImpl) Release() {
	sheets.dispatch.Release()
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
