// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"

	"github.com/go-ole/go-ole"
)

type Workbooks struct {
	dispatch *ole.IDispatch
}

func (books *Workbooks) Add() (workbook *Workbook) {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot create new excel workbook"))
			workbook = nil
		}
	}()

	workbook = new(Workbook)
	workbook.dispatch = books.call("Add")
	return workbook
}

func (books *Workbooks) Open(filePath string) (workbook *Workbook) {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot open excel file [" + filePath + "]"))
			workbook = nil
		}
	}()

	workbook = new(Workbook)
	workbook.dispatch = books.call("Open", filePath, true)
	return workbook
}

func (books *Workbooks) Close() {
	books.call("Close", false)
}

func (books *Workbooks) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(books.dispatch, methodName, parameters...)
}

func (books *Workbooks) Count() uint {
	return (uint)(getPropertyValue(books.dispatch, "Count"))
}

func (books *Workbooks) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(books.dispatch, propertyName, propertyValue)
}

func (books *Workbooks) getProperty(propertyName string) *ole.IDispatch {
	return getProperty(books.dispatch, propertyName)
}

func (books *Workbooks) getPropertyValue(propertyName string) *ole.IDispatch {
	return getProperty(books.dispatch, propertyName)
}

func (books *Workbooks) Release() {
	books.dispatch.Release()
}
