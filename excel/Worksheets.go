// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
)

type Worksheets         struct {
	dispatch *ole.IDispatch
}

func (this *Worksheets) Add() (worksheet *Worksheet) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Cannot create new excel worksheet: %s", r)
			panic(errors.New(msg))
			worksheet = nil
		}
	}()

	worksheet = new(Worksheet)
	worksheet.dispatch = this.call("Add")
	this.moveToLast(worksheet)
	return worksheet
}

func (this *Worksheets) moveToLast(worksheet *Worksheet) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Cannot excel worksheet [%s] to last position", r)
			panic(errors.New(msg))
		}
	}()

	worksheet.call("Move", nil, this.last().dispatch)
}

func (this *Worksheets) last() (worksheet *Worksheet) {
	worksheetCount := this.Count()
	lastWorksheet := this.Item(worksheetCount)
	return lastWorksheet
}

func (this *Worksheets) AddFromCsvFile(csvFilePath string, worksheetName string) (worksheet *Worksheet) {
	newWorksheet := this.Add()

	topLeftOfWorksheet := newWorksheet.Cells(1,1)

	queryTableDispatch := newWorksheet.getProperty("QueryTables" )
	newQueryTable := callMethod(queryTableDispatch, "Add", "TEXT;" + csvFilePath, topLeftOfWorksheet.dispatch)
	setProperty(newQueryTable, "TextFileParseType", 1)  // xlDelimited
	setProperty(newQueryTable, "TextFileCommaDelimiter", true)
	setProperty(newQueryTable, "TextFileSpaceDelimiter", false)
	setProperty(newQueryTable, "Refresh", false)
	newWorksheet.SetName(worksheetName)

	return newWorksheet
}

func (this *Worksheets) Count() uint {
	return (uint)(this.getPropertyValue("Count"))
}

func (this *Worksheets) Item(index uint) *Worksheet {
	worksheet := new(Worksheet)
	worksheet.dispatch = this.getProperty("Item", index)
	return worksheet
}

func (this *Worksheets) getProperty(propertyName string, parameters... interface{}) *ole.IDispatch {
	return getProperty(this.dispatch, propertyName, parameters...)
}

func (this *Worksheets) getPropertyValue(propertyName string, parameters... interface{}) int64 {
	return getPropertyValue(this.dispatch, propertyName, parameters...)
}

func (this *Worksheets) call(methodName string, parameters... interface{}) *ole.IDispatch {
	return callMethod(this.dispatch, methodName, parameters...)
}
