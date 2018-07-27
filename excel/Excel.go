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

func (handler *ExcelHandler) Initialise() error {
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	excelAppObject, err := oleutil.CreateObject("Excel.Application")

	if err != nil {
		return err
	}

	handler.appObject = excelAppObject

	newExcelIDispatch, err := excelAppObject.QueryInterface(ole.IID_IDispatch)

	if err != nil {
		return err
	}

	newExcel := new(Excel)
	newExcel.dispatch = newExcelIDispatch

	handler.excel = newExcel

	return nil
}

func InitialiseHandler() *ExcelHandler {
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

func (handler *ExcelHandler) Destroy() {
	handler.Close()
	handler.Quit()
	defer (*ole.IDispatch)(handler.excel.dispatch).Release()
	ole.CoUninitialize()
}

func (handler *ExcelHandler) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(handler.excel.dispatch, propertyName, propertyValue)
}

func (handler *ExcelHandler) getProperty(propertyName string) *ole.IDispatch {
	return getProperty(handler.excel.dispatch, propertyName)
}

func (handler *ExcelHandler) Workbooks() *Workbooks {
	if handler.workbooks == nil {
		newWorkbooks := new(Workbooks)
		newWorkbooks.dispatch = handler.getProperty("Workbooks")

		handler.workbooks = newWorkbooks
	}
	return handler.workbooks
}

func (handler *ExcelHandler) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("cannot close Excel handler")
		}
	}()

	callMethod(handler.excel.dispatch, "Save")

	workbooks := handler.Workbooks()
	workbooks.Close()
	workbooks.Release()

	return nil
}

func (handler *ExcelHandler) Quit() (err error) {
	callMethod(handler.excel.dispatch, "Quit")
	return nil
}
