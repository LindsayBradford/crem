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
	row      uint
	labelCol uint
	valueCol uint
	label    string
}

var nColsCellDetail = headerCellDetail{1, 1, 2, "ncols"}
var nRowsCellDetail = headerCellDetail{2, 1, 2, "nrows"}
var xllCornerCellDetail = headerCellDetail{3, 1, 2, "xllcorner"}
var yllCornerCellDetail = headerCellDetail{4, 1, 2, "yllcorner"}
var cellSizeCellDetail = headerCellDetail{5, 1, 2, "cellsize"}
var noDataCellDetail = headerCellDetail{6, 1, 2, "NODATA_value"}

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
	defer workbook.Close(false)

	ds.loadWorkbook(workbook)
	workbook.SetProperty("Saved", true)
}

func (ds *DataSet) loadWorkbook(workbook excel.Workbook) {
	ws := workbook.Worksheets()
	for i := uint(1); i <= ws.Count(); i++ {
		ds.loadWorksheet(ws.Item(i))
	}
	ws.Release()
}

func (ds *DataSet) loadWorksheet(sheet excel.Worksheet) {
	if isAscSheet(sheet) {
		ds.loadAscWorksheet(sheet)
	} else {
		ds.loadCsvWorksheet(sheet)
	}
}

func isAscSheet(sheet excel.Worksheet) bool {
	noDataCell := sheet.Cells(noDataCellDetail.row, noDataCellDetail.valueCol-1)
	defer noDataCell.Release()

	if value, isString := noDataCell.Value().(string); isString {
		if value == noDataCellDetail.label {
			return true
		}
	}
	return false
}

func (ds *DataSet) loadAscWorksheet(sheet excel.Worksheet) {
	newAscTable := new(tables.AscTableImpl)

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
	headerCell := sheet.Cells(detail.row, detail.valueCol)
	defer headerCell.Release()

	if valueAsDecimal, isDecimal := headerCell.Value().(float64); isDecimal {
		return valueAsDecimal
	}

	panic(detail.label + " value not retrievable")
}

func buildAscCellData(table *tables.AscTableImpl, sheet excel.Worksheet) {
	table.SetName(sheet.Name())
	table.SetColumnAndRowSize(table.Header().NumCols, table.Header().NumRows)

	for col := ascColOffset; col < table.Header().NumCols+ascColOffset; col++ {
		for row := ascRowOffset; row < table.Header().NumRows+ascRowOffset; row++ {
			cell := sheet.Cells(row, col)
			table.SetCell(col-ascColOffset, row-ascRowOffset, cell.Value())
			cell.Release()
		}
	}
}

func (ds *DataSet) loadCsvWorksheet(sheet excel.Worksheet) {
	newCsvTable := new(tables.CsvTableImpl)

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

func buildCsvCellData(table tables.CsvTable, sheet excel.Worksheet) {
	table.SetName(sheet.Name())
	colCount := excel.ColumnCount(sheet)
	rowCount := excel.RowCount(sheet)

	table.SetColumnAndRowSize(colCount, rowCount-1)

	for col := csvColOffset; col < colCount+csvColOffset; col++ {
		for row := csvRowOffset; row < rowCount+csvRowOffset-1; row++ {
			cell := sheet.Cells(row, col)
			table.SetCell(col-csvColOffset, row-csvRowOffset, cell.Value())
			cell.Release()
		}
	}
}

func (ds *DataSet) SaveAs(excelFilePath string) error {
	ds.oleWrapper(func() {
		ds.saveDataSetIntoExcelFile(excelFilePath)
	})

	return nil
}

func (ds *DataSet) saveDataSetIntoExcelFile(excelFilePath string) {
	workbooks := ds.excelHandler.Workbooks()
	defer workbooks.Release()

	workbook := workbooks.Add()
	defer workbook.Close(false)
	defer workbook.Release()

	ds.storeToWorkbook(workbook)
	ds.saveWorkbookAs(workbook, excelFilePath)
}

func (ds *DataSet) storeToWorkbook(workbook excel.Workbook) {
	worksheets := workbook.Worksheets()
	defer worksheets.Release()

	for _, table := range ds.Tables() {
		ds.storeTableToWorksheets(table, worksheets)
	}

	ds.removeEmptyDefaultWorksheet(worksheets)
	excel.ActivateFirstWorksheet(workbook)
}

func (ds *DataSet) removeEmptyDefaultWorksheet(worksheets excel.Worksheets) {
	originalWorksheet := worksheets.Item(1)
	defer originalWorksheet.Release()

	originalWorksheet.Delete()
}

func (ds *DataSet) storeTableToWorksheets(table dataset.Table, worksheets excel.Worksheets) {
	newWorksheet := worksheets.Add()
	defer newWorksheet.Release()

	newWorksheet.SetName(table.Name())

	if ascTable, isAscTable := table.(tables.AscTable); isAscTable {
		ds.storeAscTableToWorksheet(ascTable, newWorksheet)
	}
	if csvTable, isCsvTable := table.(tables.CsvTable); isCsvTable {
		ds.storeCsvTableToWorksheet(csvTable, newWorksheet)
	}

	excel.MoveWorksheetToLastInWorksheets(newWorksheet, worksheets)
	excel.AutoFitColumns(newWorksheet)
}

func (ds *DataSet) storeAscTableToWorksheet(table tables.AscTable, worksheet excel.Worksheet) {
	ds.storeAscHeaderToWorksheet(table, worksheet)
	ds.storeAscTableContentToWorksheet(table, worksheet)
}

func (ds *DataSet) storeAscHeaderToWorksheet(table tables.AscTable, worksheet excel.Worksheet) {
	ds.storeAscHeaderCellToWorksheet(worksheet, table, nColsCellDetail)
	ds.storeAscHeaderCellToWorksheet(worksheet, table, nRowsCellDetail)
	ds.storeAscHeaderCellToWorksheet(worksheet, table, xllCornerCellDetail)
	ds.storeAscHeaderCellToWorksheet(worksheet, table, yllCornerCellDetail)
	ds.storeAscHeaderCellToWorksheet(worksheet, table, cellSizeCellDetail)
	ds.storeAscHeaderCellToWorksheet(worksheet, table, noDataCellDetail)
}

func (ds *DataSet) storeAscHeaderCellToWorksheet(worksheet excel.Worksheet, table tables.AscTable, header headerCellDetail) {
	headerNameCell := worksheet.Cells(header.row, header.labelCol)
	defer headerNameCell.Release()
	headerNameCell.SetValue(header.label)

	headerValueCell := worksheet.Cells(header.row, header.valueCol)
	defer headerValueCell.Release()
	headerValueCell.SetValue(fieldForHeaderCellDetail(header, table.Header()))
}

func fieldForHeaderCellDetail(detail headerCellDetail, header tables.AscHeader) interface{} {
	switch detail {
	case nColsCellDetail:
		return header.NumCols
	case nRowsCellDetail:
		return header.NumRows
	case xllCornerCellDetail:
		return header.XllCorner
	case yllCornerCellDetail:
		return header.YllCorner
	case cellSizeCellDetail:
		return header.CellSize
	case noDataCellDetail:
		return header.NoDataValue
	}
	return nil
}

func (ds *DataSet) storeAscTableContentToWorksheet(table tables.AscTable, worksheet excel.Worksheet) {
	colSize, rowSize := table.ColumnAndRowSize()

	for col := uint(0); col < colSize; col++ {
		for row := uint(0); row < rowSize; row++ {
			cell := worksheet.Cells(row+ascRowOffset, col+ascColOffset)
			cell.SetValue(table.Cell(col, row))
			cell.Release()
		}
	}
}

func (ds *DataSet) storeCsvTableToWorksheet(table tables.CsvTable, worksheet excel.Worksheet) {
	colCount, rowCount := table.ColumnAndRowSize()

	const worksheetHeaderRow = 1

	header := table.Header()
	for col := uint(0); col < colCount; col++ {
		cell := worksheet.Cells(worksheetHeaderRow, col+csvColOffset)
		cell.SetValue(header[col])
		cell.Release()
	}

	for col := uint(0); col < colCount; col++ {
		for row := uint(0); row < rowCount; row++ {
			cell := worksheet.Cells(row+csvRowOffset, col+csvColOffset)

			value := table.Cell(col, row)
			derivedNumberFormat := deriveExcelNumberFormat(value)

			cell.SetNumberFormat(derivedNumberFormat)
			cell.SetValue(value)
			cell.Release()
		}
	}
}

func (ds *DataSet) saveWorkbookAs(workbook excel.Workbook, filePath string) {
	workbook.SaveAs(filePath)
	workbook.SetProperty("Saved", true)
}

func (ds *DataSet) Teardown() {
	ds.excelHandler.Destroy()
}

func deriveExcelNumberFormat(value interface{}) string {
	// https://www.excelhowto.com/macros/formatting-a-range-of-cells-in-excel-vba/
	switch value.(type) {
	case bool:
		return "@"
	case int, int64, uint64:
		return "#,##0"
	case float64:
		return "#,##0.00#"
	default:
		return "@"
	}
}
