// Copyright (c) 2018 Australian Rivers Institute.

package excel

import "github.com/go-ole/go-ole"

type QueryTables interface {
	AddCsvFileToWorksheet(csvFilePath string, worksheet Worksheet) QueryTable
}

type QueryTablesImpl struct {
	dispatch *ole.IDispatch
}

func (qt *QueryTablesImpl) AddCsvFileToWorksheet(csvFilePath string, worksheet Worksheet) QueryTable {
	topLeftCellOfWorksheet := worksheet.Cells(1, 1)
	cellImpl := topLeftCellOfWorksheet.(*CellImpl)

	newQueryTable := new(QueryTableImpl)
	newQueryTable.dispatch = qt.call("Add", "TEXT;"+csvFilePath, cellImpl.dispatch)
	return newQueryTable
}

func (qt *QueryTablesImpl) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(qt.dispatch, methodName, parameters...)
}
