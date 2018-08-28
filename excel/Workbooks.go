// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"

	"github.com/go-ole/go-ole"
)

const dontUpdateLinks = false
const openReadOnly = false
const dontSaveChanges = false

type Workbooks interface {
	Add() (workbook Workbook)
	Open(filePath string) (workbook Workbook)
	Close()
	Count() uint
	Release()
}

type WorkbooksImpl struct {
	oleWrapper
}

func (books *WorkbooksImpl) WithDispatch(dispatch *ole.IDispatch) *WorkbooksImpl {
	books.dispatch = dispatch
	return books
}

func (books *WorkbooksImpl) Add() (workbook Workbook) {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot create new excel workbook"))
			workbook = nil
		}
	}()

	bookDispatch := books.call("Add")
	return new(WorkbookImpl).WithDispatch(bookDispatch)
}

func (books *WorkbooksImpl) Open(filePath string) (workbook Workbook) {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot open file [" + filePath + "]"))
			workbook = nil
		}
	}()

	bookDispatch := books.call("Open", filePath, dontUpdateLinks, openReadOnly)
	return new(WorkbookImpl).WithDispatch(bookDispatch)
}

func (books *WorkbooksImpl) Close() {
	books.call("Close", dontSaveChanges)
}

func (books *WorkbooksImpl) Count() uint {
	return (uint)(getPropertyValue(books.dispatch, "Count"))
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
