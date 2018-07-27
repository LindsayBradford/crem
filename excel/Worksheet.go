// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
)

type Worksheet struct {
	dispatch *ole.IDispatch
}

func (ws *Worksheet) Name() string {
	return ws.getPropertyString("Name")
}

func (ws *Worksheet) SetName(name string) {
	ws.setProperty("Name", name)
}

func (ws *Worksheet) Delete() {
	ws.call("Delete")
}

func (ws *Worksheet) UsedRange() *Range {
	usedRange := new(Range)
	usedRange.dispatch = ws.getProperty("UsedRange")
	return usedRange
}

func (ws *Worksheet) Cells(rowIndex uint, columnIndex uint) (cell *Cell) {
	defer func() {
		if r := recover(); r != nil {
			cell = nil
		}
	}()

	cell = new(Cell)
	cell.dispatch = ws.getProperty("Cells", rowIndex, columnIndex)
	return cell
}

func (ws *Worksheet) getProperty(propertyName string, parameters ...interface{}) *ole.IDispatch {
	return getProperty(ws.dispatch, propertyName, parameters...)
}

func (ws *Worksheet) getPropertyString(propertyName string, parameters ...interface{}) string {
	return getPropertyString(ws.dispatch, propertyName, parameters...)
}

func (ws *Worksheet) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(ws.dispatch, propertyName, propertyValue)
}

func (ws *Worksheet) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(ws.dispatch, methodName, parameters...)
}
