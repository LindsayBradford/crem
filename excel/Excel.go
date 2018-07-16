// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
"errors"


"github.com/go-ole/go-ole"
"github.com/go-ole/go-ole/oleutil"

)

type ApplicationObject *ole.IUnknown

type Excel struct {
	dispatch *ole.IDispatch
}


type ExcelHandler struct {
	appObject ApplicationObject
	excel     *Excel
	workbooks *Workbooks
}

func (this *ExcelHandler) Initialise() error {
	excelAppObject, err := oleutil.CreateObject("Excel.Application")

	if (err != nil) {
		return err
	}

	this.appObject = excelAppObject

	newExcelIDispatch, err := excelAppObject.QueryInterface(ole.IID_IDispatch)

	if (err != nil) {
		return err
	}

	newExcel := new(Excel)
	newExcel.dispatch = newExcelIDispatch

	this.excel = newExcel

	return nil
}

func InitialiseHandler() (*ExcelHandler) {
	ole.CoInitialize(0)

	newHandler := new(ExcelHandler)

	initialiseError := newHandler.Initialise()
	if initialiseError != nil {
		panic(initialiseError)
	}

	newHandler.setProperty("Visible", false)
	newHandler.setProperty("DisplayAlerts", false)
	newHandler.setProperty("ScreenUpdating", false)

	return newHandler
}

func (this *ExcelHandler) Destroy() {
	this.Close()
	this.Quit()
	defer (*ole.IDispatch)(this.excel.dispatch).Release()
	ole.CoUninitialize()
}

func (this *ExcelHandler) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(this.excel.dispatch, propertyName, propertyValue)
}

func (this *ExcelHandler) getProperty(propertyName string) *ole.IDispatch {
	return getProperty(this.excel.dispatch, propertyName)
}

func (this *ExcelHandler) Workbooks() *Workbooks {
	if this.workbooks == nil {
		newWorkbooks := new(Workbooks)
		newWorkbooks.dispatch = this.getProperty("Workbooks")

		this.workbooks = newWorkbooks
	}
	return this.workbooks
}

func (this *ExcelHandler) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Cannot save Excel")
		}
	}()

	callMethod(this.excel.dispatch, "Save")

	workbooks := this.Workbooks()
	workbooks.Close(); workbooks.Release()

	return nil
}

func (this *ExcelHandler) Quit() (err error) {
	callMethod(this.excel.dispatch, "Quit")
	return nil
}