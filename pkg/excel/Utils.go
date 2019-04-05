// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func getProperty(iDispatch *ole.IDispatch, propertyName string, parameters ...interface{}) *ole.IDispatch {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).ToIDispatch()
}

func getPropertyValue(iDispatch *ole.IDispatch, propertyName string, parameters ...interface{}) int64 {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).Val
}

func getPropertyVariant(iDispatch *ole.IDispatch, propertyName string, parameters ...interface{}) interface{} {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).Value()
}

func getPropertyString(iDispatch *ole.IDispatch, propertyName string, parameters ...interface{}) string {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).ToString()
}

func setProperty(iDispatch *ole.IDispatch, propertyName string, propertyValue interface{}) {
	oleutil.PutProperty(iDispatch, propertyName, propertyValue)
}

func callMethod(dispatch *ole.IDispatch, methodName string, parameters ...interface{}) *ole.IDispatch {
	return oleutil.MustCallMethod(dispatch, methodName, parameters...).ToIDispatch()
}

func AddWorksheetFromCsvFileToWorkbook(csvFilePath string, worksheetName string, workbook Workbook) Worksheet {
	worksheets := workbook.Worksheets()
	defer worksheets.Release()
	newWorksheet := worksheets.Add()

	AddCsvFileContentToWorksheet(csvFilePath, newWorksheet)

	newWorksheet.SetName(worksheetName)
	MoveWorksheetToLastInWorksheets(newWorksheet, worksheets)

	return newWorksheet
}

func AddCsvFileContentToWorksheet(csvFilePath string, worksheet Worksheet) {
	ClearUsedRange(worksheet)

	queryTables := worksheet.QueryTables()
	defer queryTables.Release()

	newQueryTable := queryTables.AddCsvFileToWorksheet(csvFilePath, worksheet)
	newQueryTable.SetProperty("TextFileParseType", 1) // xlDelimited
	newQueryTable.SetProperty("TextFileCommaDelimiter", true)
	newQueryTable.SetProperty("TextFileSpaceDelimiter", false)
	newQueryTable.SetProperty("Refresh", false)

	newQueryTable.Release()
}

func LastOfWorksheets(worksheets Worksheets) (worksheet Worksheet) {
	worksheetCount := worksheets.Count()
	return worksheets.Item(worksheetCount)
}

func ColumnCount(worksheet Worksheet) uint {
	usedRange := worksheet.UsedRange()
	defer usedRange.Release()
	columns := usedRange.Columns()
	defer columns.Release()
	return columns.Count()
}

func RowCount(worksheet Worksheet) uint {
	usedRange := worksheet.UsedRange()
	defer usedRange.Release()
	rows := usedRange.Rows()
	defer rows.Release()
	return rows.Count()
}

func MoveWorksheetToLastInWorksheets(worksheet Worksheet, worksheets Worksheets) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("cannot move excel worksheet [%s] to last position", r)
			panic(errors.New(msg))
		}
	}()
	lastWorksheet := LastOfWorksheets(worksheets)
	defer lastWorksheet.Release()
	worksheet.MoveToAfterWorksheet(lastWorksheet)
}

func AutoFitColumns(worksheet Worksheet) {
	usedRange := worksheet.UsedRange()
	defer usedRange.Release()
	columns := usedRange.Columns()
	defer columns.Release()
	columns.AutoFit()
}

func ClearUsedRange(worksheet Worksheet) {
	usedRange := worksheet.UsedRange()
	defer usedRange.Release()

	usedRange.Clear()
}

func ActivateFirstWorksheet(workbook Workbook) {
	worksheets := workbook.Worksheets()
	defer worksheets.Release()

	firstWorksheet := worksheets.Item(1)
	defer firstWorksheet.Release()

	firstWorksheet.Activate()
}
