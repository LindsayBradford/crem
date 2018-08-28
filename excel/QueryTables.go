// Copyright (c) 2018 Australian Rivers Institute.

package excel

import "github.com/go-ole/go-ole"

type QueryTables interface {
	AddCsvFileToWorksheet(csvFilePath string, worksheet Worksheet) QueryTable
	Release()
}

type QueryTablesImpl struct {
	oleWrapper
}

func (qt *QueryTablesImpl) WithDispatch(dispatch *ole.IDispatch) *QueryTablesImpl {
	qt.dispatch = dispatch
	return qt
}

func (qt *QueryTablesImpl) AddCsvFileToWorksheet(csvFilePath string, worksheet Worksheet) QueryTable {
	topLeftCellOfWorksheet := worksheet.Cells(1, 1)
	defer topLeftCellOfWorksheet.Release()
	cellImpl := topLeftCellOfWorksheet.(*CellImpl)
	dispatch := qt.call("Add", "TEXT;"+csvFilePath, cellImpl.dispatch)
	return new(QueryTableImpl).WithDispatch(dispatch)
}

func (qt *QueryTablesImpl) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(qt.dispatch, methodName, parameters...)
}
