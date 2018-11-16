// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
)

type Range interface {
	Rows() Range
	Columns() Range
	Count() uint
	Clear()
	AutoFit()
	Release()
}

type RangeImpl struct {
	dispatch *ole.IDispatch
}

func (r *RangeImpl) WithDispatch(dispatch *ole.IDispatch) *RangeImpl {
	r.dispatch = dispatch
	return r
}

func (r *RangeImpl) Rows() Range {
	dispatch := r.getProperty("Rows")
	return new(RangeImpl).WithDispatch(dispatch)
}

func (r *RangeImpl) Columns() Range {
	dispatch := r.getProperty("Columns")
	return new(RangeImpl).WithDispatch(dispatch)
}

func (r *RangeImpl) Count() uint {
	return (uint)(r.getPropertyValue("Count"))
}

func (r *RangeImpl) Clear() {
	r.call("Clear")
}

func (r *RangeImpl) AutoFit() {
	r.call("AutoFit")
}

func (r *RangeImpl) Release() {
	r.dispatch.Release()
}

func (r *RangeImpl) getProperty(propertyName string, parameters ...interface{}) *ole.IDispatch {
	return getProperty(r.dispatch, propertyName, parameters...)
}

func (r *RangeImpl) getPropertyValue(propertyName string, parameters ...interface{}) int64 {
	return getPropertyValue(r.dispatch, propertyName, parameters...)
}

func (r *RangeImpl) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(r.dispatch, methodName, parameters...)
}
