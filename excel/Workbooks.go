// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"

	"github.com/go-ole/go-ole"
)

type Workbooks interface {
	Add() (workbook Workbook)
	Open(filePath string) (workbook Workbook)
	Close()
	Count() uint
	Release()
}

type WorkbooksImpl struct {
	dispatch *ole.IDispatch
}

func (books *WorkbooksImpl) Add() (workbook Workbook) {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot create new excel workbook"))
			workbook = nil
		}
	}()

	newWorkbook := new(WorkbookImpl)
	newWorkbook.dispatch = books.call("Add")
	return newWorkbook
}

func (books *WorkbooksImpl) Open(filePath string) (workbook Workbook) {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot open file [" + filePath + "]"))
			workbook = nil
		}
	}()

	newWorkbook := new(WorkbookImpl)
	newWorkbook.dispatch = books.call("Open", filePath, true)
	return newWorkbook
}

func (books *WorkbooksImpl) Close() {
	books.call("Close", false)
}

func (books *WorkbooksImpl) Count() uint {
	return (uint)(getPropertyValue(books.dispatch, "Count"))
}

func (books *WorkbooksImpl) Release() {
	books.dispatch.Release()
}

func (books *WorkbooksImpl) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(books.dispatch, methodName, parameters...)
}

func (books *WorkbooksImpl) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(books.dispatch, propertyName, propertyValue)
}

func (books *WorkbooksImpl) getProperty(propertyName string) *ole.IDispatch {
	return getProperty(books.dispatch, propertyName)
}

func (books *WorkbooksImpl) getPropertyValue(propertyName string) *ole.IDispatch {
	return getProperty(books.dispatch, propertyName)
}
