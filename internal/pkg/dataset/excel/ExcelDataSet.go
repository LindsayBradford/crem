// Copyright (c) 2018 Australian Rivers Institute.

package excel

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/pkg/excel"
	"github.com/LindsayBradford/crem/pkg/threading"
)

func NewDataSet(name string, oleWrapper threading.MainThreadFunctionWrapper) *DataSet {
	excelDataSet := new(DataSet).WithOleFunctionWrapper(oleWrapper)
	excelDataSet.excelHandler.Initialise()
	excelDataSet.Initialise(name)
	return excelDataSet
}

type headerCellDetail struct {
	row   uint
	col   uint
	label string
}

var nColsCellDetail = headerCellDetail{1, 2, "ncols"}
var nRowsCellDetail = headerCellDetail{2, 2, "nrows"}
var xllCornerCellDetail = headerCellDetail{3, 2, "xllcorner"}
var yllCornerCellDetail = headerCellDetail{4, 2, "yllcorner"}
var cellSizeCellDetail = headerCellDetail{5, 2, "cellsize"}
var noDataCellDetail = headerCellDetail{6, 2, "NODATA_value"}

const ascRowOffset = uint(7)
const ascColOffset = uint(1)

const csvRowOffset = uint(2)
const csvColOffset = uint(1)

type DataSet struct {
	dataset.DataSetImpl

	excelHandler     excel.Handler
	absoluteFilePath string
	oleWrapper       threading.MainThreadFunctionWrapper
}

func (ds *DataSet) WithOleFunctionWrapper(wrapper threading.MainThreadFunctionWrapper) *DataSet {
	ds.oleWrapper = wrapper
	return ds
}

func (ds *DataSet) Load(excelFilePath string) error {
	ds.oleWrapper(func() {
		ds.loadExcelFileIntoDataSet(excelFilePath)
	})

	return nil
}

func (ds *DataSet) loadExcelFileIntoDataSet(excelFilePath string) {
	workbooks := ds.excelHandler.Workbooks()
	defer workbooks.Release()

	workbook := workbooks.Open(excelFilePath)
	ds.processWorkbook(workbook)

	workbook.SetProperty("Saved", true)
	workbook.Close(false)
}

func (ds *DataSet) processWorkbook(workbook excel.Workbook) {
	ws := workbook.Worksheets()
	for i := uint(1); i <= ws.Count(); i++ {
		ds.processWorksheet(ws.Item(i))
	}
	ws.Release()
}

func (ds *DataSet) processWorksheet(sheet excel.Worksheet) {
	if isAscSheet(sheet) {
		ds.processAscWorksheet(sheet)
	} else {
		ds.processCsvWorksheet(sheet)
	}
}

func isAscSheet(sheet excel.Worksheet) bool {
	noDataCell := sheet.Cells(noDataCellDetail.row, noDataCellDetail.col-1)
	defer noDataCell.Release()

	if value, isString := noDataCell.Value().(string); isString {
		if value == noDataCellDetail.label {
			return true
		}
	}
	return false
}

func (ds *DataSet) processAscWorksheet(sheet excel.Worksheet) {
	newAscTable := new(tables.AscTable)

	newAScHeader := buildAscHeader(sheet)
	newAscTable.SetHeader(newAScHeader)
	buildAscCellData(newAscTable, sheet)

	ds.AddTable(sheet.Name(), newAscTable)
}

func buildAscHeader(sheet excel.Worksheet) tables.AscHeader {
	newAScHeader := tables.AscHeader{}

	newAScHeader.NumCols = uint(retrieveHeaderValue(nColsCellDetail, sheet))
	newAScHeader.NumRows = uint(retrieveHeaderValue(nRowsCellDetail, sheet))
	newAScHeader.XllCorner = retrieveHeaderValue(xllCornerCellDetail, sheet)
	newAScHeader.YllCorner = retrieveHeaderValue(yllCornerCellDetail, sheet)
	newAScHeader.CellSize = int64(retrieveHeaderValue(cellSizeCellDetail, sheet))
	newAScHeader.NoDataValue = int64(retrieveHeaderValue(noDataCellDetail, sheet))

	return newAScHeader
}

func retrieveHeaderValue(detail headerCellDetail, sheet excel.Worksheet) float64 {
	headerCell := sheet.Cells(detail.row, detail.col)
	defer headerCell.Release()

	if valueAsDecimal, isDecimal := headerCell.Value().(float64); isDecimal {
		return valueAsDecimal
	}

	panic(detail.label + " value not retrievable")
}

func buildAscCellData(table *tables.AscTable, sheet excel.Worksheet) {
	table.SetName(sheet.Name())
	table.SetSize(table.Header().NumCols, table.Header().NumRows)

	for col := ascColOffset; col < table.Header().NumCols+ascColOffset; col++ {
		for row := ascRowOffset; row < table.Header().NumRows+ascRowOffset; row++ {
			cell := sheet.Cells(row, col)
			table.SetCell(col-ascColOffset, row-ascRowOffset, cell.Value())
			cell.Release()
		}
	}
}

func (ds *DataSet) processCsvWorksheet(sheet excel.Worksheet) {
	newCsvTable := new(tables.CsvTable)

	newCsvHeader := buildCsvHeader(sheet)
	newCsvTable.SetHeader(newCsvHeader)
	buildCsvCellData(newCsvTable, sheet)

	ds.AddTable(sheet.Name(), newCsvTable)
}

func buildCsvHeader(sheet excel.Worksheet) tables.CsvHeader {
	newCsvHeader := make(tables.CsvHeader, 0)

	colCount := excel.ColumnCount(sheet)
	for col := uint(1); col <= colCount; col++ {
		headerCell := sheet.Cells(1, col)
		defer headerCell.Release()

		headerValue := headerCell.Value().(string)
		newCsvHeader = append(newCsvHeader, headerValue)
	}

	return newCsvHeader
}

func buildCsvCellData(table *tables.CsvTable, sheet excel.Worksheet) {
	table.SetName(sheet.Name())
	colCount := excel.ColumnCount(sheet)
	rowCount := excel.RowCount(sheet)

	table.SetSize(colCount, rowCount-1)

	for col := csvColOffset; col < colCount+csvColOffset; col++ {
		for row := csvRowOffset; row < rowCount+csvRowOffset-1; row++ {
			cell := sheet.Cells(row, col)
			table.SetCell(col-csvColOffset, row-csvRowOffset, cell.Value())
			cell.Release()
		}
	}
}
