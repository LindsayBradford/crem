// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"

	"github.com/go-ole/go-ole"
)

type Workbooks         struct {
	dispatch *ole.IDispatch
}

func (this *Workbooks) Add() (workbook *Workbook, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Cannot create new excel workbook")
			workbook = nil
		}
	}()

	workbook = new(Workbook)
	workbook.dispatch = this.call("Add")
	return workbook, nil
}

func (this *Workbooks) Open(filePath string) (workbook *Workbook, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Cannot open excel file [" + filePath + "]")
			workbook = nil
		}
	}()

	workbook = new(Workbook)
	workbook.dispatch = this.call("Open", filePath, true)
	return workbook, nil
}

func (this *Workbooks) Close() {
	this.call("Close", false)
}

func (this *Workbooks) call(methodName string, parameters... interface{}) *ole.IDispatch {
	return callMethod(this.dispatch, methodName, parameters...)
}

func (this *Workbooks) Count() uint {
	return	(uint)(getPropertyValue(this.dispatch, "Count"))
}

func (this *Workbooks) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(this.dispatch, propertyName, propertyValue)
}

func (this *Workbooks) getProperty(propertyName string) *ole.IDispatch {
	return getProperty(this.dispatch, propertyName)
}

func (this *Workbooks) getPropertyValue(propertyName string) *ole.IDispatch {
	return getProperty(this.dispatch, propertyName)
}

func (this *Workbooks) Release() {
	this.dispatch.Release()
}
