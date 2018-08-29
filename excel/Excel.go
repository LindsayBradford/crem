// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/pkg/errors"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type oleWrapper struct {
	dispatch *ole.IDispatch
}

func (ep *oleWrapper) Release() {
	ep.dispatch.Release()
}

type Excel struct {
	oleWrapper
}

func (excel *Excel) WithDispatch(dispatch *ole.IDispatch) *Excel {
	excel.dispatch = dispatch
	return excel
}

type Handler struct {
	excel *Excel
}

func (handler *Handler) Initialise() *Handler {
	defer func() {
		if r := recover(); r != nil {
			recoveredError, ok := r.(error)
			if ok {
				wrappedError := errors.Wrap(recoveredError, "excel handler initialise failed")
				panic(wrappedError)
			}
			panic(r)
		}
	}()

	appObject, err := oleutil.CreateObject("Excel.Application")
	if err != nil {
		panic(err)
	}
	defer appObject.Release()

	excelDispatch, err := appObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		panic(err)
	}

	handler.excel = new(Excel).WithDispatch(excelDispatch)
	handler.setPropertiesForSilentOperation()

	return handler
}

func (handler *Handler) setPropertiesForSilentOperation() {
	handler.setProperty("Visible", false)
	handler.setProperty("DisplayAlerts", false)
	handler.setProperty("ScreenUpdating", false)
}

func (handler *Handler) Destroy() {
	handler.excel.Release()
}

func (handler *Handler) Workbooks() Workbooks {
	workbooksDispatch := handler.getProperty("Workbooks")
	newWorkbooks := new(WorkbooksImpl).WithDispatch(workbooksDispatch)
	return newWorkbooks
}

func (handler *Handler) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(handler.excel.dispatch, propertyName, propertyValue)
}

func (handler *Handler) getProperty(propertyName string) *ole.IDispatch {
	return getProperty(handler.excel.dispatch, propertyName)
}
