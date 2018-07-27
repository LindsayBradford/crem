// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
)

type Range struct {
	dispatch *ole.IDispatch
}

func (r *Range) Rows() *Range {
	rows := new(Range)
	rows.dispatch = r.getProperty("Rows")
	return rows
}

func (r *Range) Columns() *Range {
	columns := new(Range)
	columns.dispatch = r.getProperty("Columns")
	return columns
}

func (r *Range) Count() uint {
	return (uint)(r.getPropertyValue("Count"))
}

func (r *Range) Clear() {
	r.call("Clear")
}

func (r *Range) AutoFit() {
	r.call("AutoFit")
}

func (r *Range) getProperty(propertyName string, parameters ...interface{}) *ole.IDispatch {
	return getProperty(r.dispatch, propertyName, parameters...)
}

func (r *Range) getPropertyValue(propertyName string, parameters ...interface{}) int64 {
	return getPropertyValue(r.dispatch, propertyName, parameters...)
}

func (r *Range) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(r.dispatch, methodName, parameters...)
}
