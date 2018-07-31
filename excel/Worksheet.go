// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
)

type Worksheet interface {
	Name() string
	SetName(name string)
	Delete()
	UsedRange() Range
	Cells(rowIndex uint, columnIndex uint) (cell Cell)
	QueryTables() QueryTables
	MoveToAfterWorksheet(worksheet Worksheet)
}

type WorksheetImpl struct {
	dispatch *ole.IDispatch
}

func (ws *WorksheetImpl) Name() string {
	return ws.getPropertyString("Name")
}

func (ws *WorksheetImpl) SetName(name string) {
	ws.setProperty("Name", name)
}

func (ws *WorksheetImpl) Delete() {
	ws.call("Delete")
}

func (ws *WorksheetImpl) UsedRange() Range {
	usedRange := new(RangeImpl)
	usedRange.dispatch = ws.getProperty("UsedRange")
	return usedRange
}

func (ws *WorksheetImpl) Cells(rowIndex uint, columnIndex uint) (cell Cell) {
	defer func() {
		if r := recover(); r != nil {
			cell = nil
		}
	}()

	newCell := new(CellImpl)
	newCell.dispatch = ws.getProperty("Cells", rowIndex, columnIndex)
	return newCell
}

func (ws *WorksheetImpl) QueryTables() QueryTables {
	newQueryTables := new(QueryTablesImpl)
	newQueryTables.dispatch = ws.getProperty("QueryTables")
	return newQueryTables
}

func (ws *WorksheetImpl) MoveToAfterWorksheet(worksheet Worksheet) {
	worksheetImpl := worksheet.(*WorksheetImpl)
	ws.call("Move", nil, worksheetImpl.dispatch)
}

func (ws *WorksheetImpl) getProperty(propertyName string, parameters ...interface{}) *ole.IDispatch {
	return getProperty(ws.dispatch, propertyName, parameters...)
}

func (ws *WorksheetImpl) getPropertyString(propertyName string, parameters ...interface{}) string {
	return getPropertyString(ws.dispatch, propertyName, parameters...)
}

func (ws *WorksheetImpl) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(ws.dispatch, propertyName, propertyValue)
}

func (ws *WorksheetImpl) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(ws.dispatch, methodName, parameters...)
}
