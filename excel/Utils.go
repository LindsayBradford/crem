// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func getProperty(iDispatch *ole.IDispatch, propertyName string, parameters... interface{}) *ole.IDispatch {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).ToIDispatch()
}

func getPropertyValue(iDispatch *ole.IDispatch, propertyName string, parameters... interface{}) int64 {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).Val
}

func getPropertyVariant(iDispatch *ole.IDispatch, propertyName string, parameters... interface{}) interface{} {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).Value()
}

func getPropertyString(iDispatch *ole.IDispatch, propertyName string, parameters... interface{}) string {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).ToString()
}

func setProperty(iDispatch *ole.IDispatch, propertyName string, propertyValue interface{}) {
	oleutil.PutProperty(iDispatch, propertyName, propertyValue)
}

func callMethod(dispatch *ole.IDispatch, methodName string, parameters... interface{}) *ole.IDispatch {
	return oleutil.MustCallMethod(dispatch, methodName, parameters...).ToIDispatch()
}

func AddWorksheetFromCsvFileToWorkbook(csvFilePath string, worksheetName string, workbook *Workbook) *Worksheet {
	worksheets := workbook.Worksheets()
	newWorksheet :=  worksheets.Add()

	AddCsvFileContentToWorksheet(csvFilePath, newWorksheet)

	newWorksheet.SetName(worksheetName)
	MoveWorksheetToLastInWorksheets(newWorksheet, worksheets)

	return newWorksheet
}

func AddCsvFileContentToWorksheet(csvFilePath string, worksheet *Worksheet) {
	worksheet.UsedRange().Clear()
	topLeftOfWorksheet := worksheet.Cells(1,1)

	queryTableDispatch := worksheet.getProperty("QueryTables" )
	newQueryTable := callMethod(queryTableDispatch, "Add", "TEXT;" + csvFilePath, topLeftOfWorksheet.dispatch)
	setProperty(newQueryTable, "TextFileParseType", 1)  // xlDelimited
	setProperty(newQueryTable, "TextFileCommaDelimiter", true)
	setProperty(newQueryTable, "TextFileSpaceDelimiter", false)
	setProperty(newQueryTable, "Refresh", false)
}

func LastOfWorksheets(worksheets *Worksheets) (worksheet *Worksheet) {
	worksheetCount := worksheets.Count()
	lastWorksheet := worksheets.Item(worksheetCount)
	return lastWorksheet
}

func MoveWorksheetToLastInWorksheets(worksheet *Worksheet, worksheets *Worksheets) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Cannot excel worksheet [%s] to last position", r)
			panic(errors.New(msg))
		}
	}()

	worksheet.call("Move", nil, LastOfWorksheets(worksheets).dispatch)
}