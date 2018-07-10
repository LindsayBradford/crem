// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"

	"github.com/go-ole/go-ole"
)

type Worksheets         struct {
	dispatch *ole.IDispatch
}

func (this *Worksheets) Add() (worksheet *Worksheet, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Cannot create new excel worksheet")
			worksheet = nil
		}
	}()

	worksheet = new(Worksheet)
	worksheet.dispatch = this.call("Add")
	return worksheet, nil
}

func (this *Worksheets) Count() uint {
	return (uint)(this.getPropertyValue("Count"))
}

func (this *Worksheets) getPropertyValue(propertyName string, parameters... interface{}) int64 {
	return getPropertyValue(this.dispatch, propertyName, parameters...)
}

func (this *Worksheets) call(methodName string, parameters... interface{}) *ole.IDispatch {
	return callMethod(this.dispatch, methodName, parameters...)
}
