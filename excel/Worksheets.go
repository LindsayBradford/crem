// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
)

type Worksheets struct {
	dispatch *ole.IDispatch
}

func (sheets *Worksheets) Add() (worksheet *Worksheet) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Cannot create new excel worksheet: %s", r)
			panic(errors.New(msg))
			worksheet = nil
		}
	}()

	worksheet = new(Worksheet)
	worksheet.dispatch = sheets.call("Add")
	return worksheet
}

func (sheets *Worksheets) Count() uint {
	return (uint)(sheets.getPropertyValue("Count"))
}

func (sheets *Worksheets) Item(index uint) *Worksheet {
	worksheet := new(Worksheet)
	worksheet.dispatch = sheets.getProperty("Item", index)
	return worksheet
}

func (sheets *Worksheets) getProperty(propertyName string, parameters ...interface{}) *ole.IDispatch {
	return getProperty(sheets.dispatch, propertyName, parameters...)
}

func (sheets *Worksheets) getPropertyValue(propertyName string, parameters ...interface{}) int64 {
	return getPropertyValue(sheets.dispatch, propertyName, parameters...)
}

func (sheets *Worksheets) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(sheets.dispatch, methodName, parameters...)
}
