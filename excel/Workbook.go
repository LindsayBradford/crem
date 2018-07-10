// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
)

type Workbook         struct {
	dispatch *ole.IDispatch
}

func (this *Workbook) Worksheets() (worksheets *Worksheets, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Cannot retrieve worksheets")
			worksheets = nil
		}
	}()

	worksheets = new(Worksheets)
	worksheets.dispatch = this.getProperty("Worksheets")
	return worksheets, nil
}

func (this *Workbook) getProperty(propertyName string, parameters ... interface{}) *ole.IDispatch {
	return getProperty(this.dispatch, propertyName, parameters...)
}

func (this *Workbook) Worksheet(index uint) (worksheet *Worksheet, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Cannot open worksheet at index [" + fmt.Sprintf("%d", index) + "]")
			worksheet = nil
		}
	}()

	worksheet = new(Worksheet)
	worksheet.dispatch = this.getProperty("Worksheets", index)
	return worksheet, nil
}

func (this *Workbook) Close() {
	this.call("Close")
}

func (this *Workbook) call(methodName string, parameters... interface{}) *ole.IDispatch {
	return callMethod(this.dispatch, methodName, parameters...)
}

