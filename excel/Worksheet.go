// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
)

type Worksheet         struct {
	dispatch *ole.IDispatch
}


func (this *Worksheet) Name() string {
	return this.getPropertyString("Name")
}

func (this *Worksheet) SetName(name string) {
	this.setProperty("Name", name)
}

func (this *Worksheet) UsedRange() (usedRange *Range) {
	usedRange = new(Range)
	usedRange.dispatch = this.getProperty("UsedRange")
	return usedRange
}

func (this *Worksheet) Cells(rowIndex uint , columnIndex uint ) (cell *Cell) {
	defer func() {
		if r := recover(); r != nil {
			cell = nil
		}
	}()

	newCell := new(Cell)
	newCell.dispatch = this.getProperty("Cells", rowIndex, columnIndex)
	return newCell
}

func (this *Worksheet) getProperty(propertyName string, parameters... interface{})  *ole.IDispatch {
	return getProperty(this.dispatch, propertyName, parameters...)
}

func (this *Worksheet) getPropertyString(propertyName string, parameters... interface{}) string {
	return getPropertyString(this.dispatch, propertyName, parameters...)
}

func (this *Worksheet) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(this.dispatch, propertyName, propertyValue)
}

func (this *Worksheet) call(methodName string, parameters... interface{}) *ole.IDispatch {
	return callMethod(this.dispatch, methodName, parameters...)
}
