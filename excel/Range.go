// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
)

type Range struct {
	dispatch *ole.IDispatch
}

func (this *Range) Rows() (rows *Range) {
	rows = new(Range)
	rows.dispatch = this.getProperty("Rows")
	return rows
}

func (this *Range) Columns() (columns *Range) {
	columns = new(Range)
	columns.dispatch = this.getProperty("Columns")
	return columns
}

func (this *Range) Count() uint {
	return (uint)(this.getPropertyValue("Count"))
}

func (this *Range) Clear() {
	this.call("Clear")
}

func (this *Range) AutoFit() {
	this.call("AutoFit")
}

func (this *Range) getProperty(propertyName string, parameters... interface{})  *ole.IDispatch {
	return getProperty(this.dispatch, propertyName, parameters...)
}

func (this *Range) getPropertyValue(propertyName string, parameters... interface{}) int64 {
	return getPropertyValue(this.dispatch, propertyName, parameters...)
}

func (this *Range) call(methodName string, parameters... interface{}) *ole.IDispatch {
	return callMethod(this.dispatch, methodName, parameters...)
}
