// Copyright (c) 2019 Australian Rivers Institute.

package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	myErrors "github.com/LindsayBradford/crem/pkg/errors"
	myStrings "github.com/LindsayBradford/crem/pkg/strings"
	"github.com/pkg/errors"
)

var caster *myStrings.BaseCaster

func init() {
	caster = new(myStrings.BaseCaster).WithNumbersAsFloats()
}

func NewDataSet(name string) *DataSet {
	dataSet := new(DataSet)
	dataSet.Initialise(name)

	dataSet.errors = myErrors.New("Csv Dataset Errors")
	return dataSet
}

var metaTableHeadings = [2]string{"TableName", "FilePath"}

type DataSet struct {
	dataset.DataSetImpl
	filePath string
	errors   *myErrors.CompositeError
}

func (ds *DataSet) Load(baseCsvFilePath string) error {
	ds.filePath = baseCsvFilePath
	pathInfo, err := os.Stat(baseCsvFilePath)
	if os.IsNotExist(err) {
		newError := errors.Errorf("file specified [%s] does not exist", baseCsvFilePath)
		ds.errors.Add(newError)
	}
	if pathInfo != nil && pathInfo.Mode().IsDir() {
		newError := errors.Errorf("file specified [%s] is a directory, not a file", baseCsvFilePath)
		ds.errors.Add(newError)
	}

	ds.loadMetaFile(baseCsvFilePath)

	if ds.errors.Size() > 0 {
		return ds.errors
	}
	return nil
}

func (ds *DataSet) Errors() error {
	if ds.errors.Size() > 0 {
		return ds.errors
	}
	return nil
}

func (ds *DataSet) loadMetaFile(metaCsvFilePath string) {
	metaTable := ds.loadCsvIntoTable(metaCsvFilePath)
	if ds.errors.Size() > 0 {
		return
	}

	ds.verifyMetaTable(metaTable)
	if ds.errors.Size() > 0 {
		return
	}

	metaTableName := ds.deriveMetaTableName(metaCsvFilePath)
	ds.AddTable(metaTableName, metaTable)
}

func (ds *DataSet) verifyMetaTable(metaTable tables.CsvTable) {
	ds.verifyMetaTableHeader(metaTable)
	ds.verifyMetaTableRows(metaTable)
}

func (ds *DataSet) verifyMetaTableHeader(metaTable tables.CsvTable) {
	header := metaTable.Header()

	actualColSize, _ := metaTable.ColumnAndRowSize()
	expectedColSize := uint(len(metaTableHeadings))

	if actualColSize != expectedColSize {
		errorMsg := fmt.Sprintf("expected csv meta-table heading columns to equal [%d], but got [%d]", expectedColSize, actualColSize)
		newError := errors.New(errorMsg)
		ds.errors.Add(newError)
	}

	for headingIndex := range metaTableHeadings {
		expectedHeading := metaTableHeadings[headingIndex]
		actualHeading := header[headingIndex]
		if actualHeading != expectedHeading {
			newError := errors.New("expected csv meta-table heading [" + expectedHeading + "], but got [" + actualHeading + "]")
			ds.errors.Add(newError)
		}
	}
}

func (ds *DataSet) verifyMetaTableRows(metaTable tables.CsvTable) {
	dataSetPath := filepath.Dir(ds.filePath)
	oldWorkingDirectory, _ := os.Getwd()

	os.Chdir(dataSetPath)
	defer os.Chdir(oldWorkingDirectory)

	_, rowSize := metaTable.ColumnAndRowSize()
	const tableNameCol = 0
	const filePathCol = 1

	preValidationErrorCount := ds.errors.Size()

	for rowIndex := uint(0); rowIndex < rowSize; rowIndex++ {
		filePathAtRow := metaTable.Cell(filePathCol, rowIndex).(string)

		pathInfo, err := os.Stat(filePathAtRow)
		if os.IsNotExist(err) {
			newError := errors.Errorf("meta-file at row [%d], file specified [%s] does not exist", rowIndex, filePathAtRow)
			ds.errors.Add(newError)
		}
		if pathInfo != nil && pathInfo.Mode().IsDir() {
			newError := errors.Errorf("meta-file at row [%d], file specified [%s] is a directory, not a file", rowIndex, filePathAtRow)
			ds.errors.Add(newError)
		}

		tableNameAtRow := metaTable.Cell(tableNameCol, rowIndex).(string)
		loadedTableAtRow := ds.loadCsvIntoTable(filePathAtRow)

		if preValidationErrorCount == ds.errors.Size() {
			ds.AddTable(tableNameAtRow, loadedTableAtRow)
		}
	}
}

func (ds *DataSet) loadCsvIntoTable(csvFilePath string) tables.CsvTable {
	records, loadError := loadCsvRecords(csvFilePath)
	if loadError != nil {
		ds.errors.Add(loadError)
		return nil
	}

	return ds.deriveTableFromRecords(records)
}

func (ds *DataSet) ParseCsvTextIntoTable(tableName string, csvContent string) {
	records, parseError := parseCsvText(csvContent)
	if parseError != nil {
		ds.errors.Add(parseError)
		return
	}

	derivedTable := ds.deriveTableFromRecords(records)
	ds.AddTable(tableName, derivedTable)
}

type tableContext struct {
	records [][]string
	rowSize uint
	colSize uint

	table *tables.CsvTableImpl
}

func (ds *DataSet) deriveTableFromRecords(inputRecords [][]string) tables.CsvTable {
	context := ds.deriveContextFromRecords(inputRecords)

	ds.assignTableHeaders(context)
	ds.assignTableContent(context)

	return context.table
}

func (ds *DataSet) assignTableContent(context *tableContext) {
	for rowIndex := uint(1); rowIndex <= context.rowSize; rowIndex++ {
		for colIndex := uint(0); colIndex < context.colSize; colIndex++ {
			recordAsBaseType := toBaseType(context.records[rowIndex][colIndex])
			context.table.SetCell(colIndex, rowIndex-1, recordAsBaseType)
		}
	}
}

func (ds *DataSet) assignTableHeaders(context *tableContext) {
	newCsvHeader := deriveHeader(context.records)
	context.table.SetHeader(newCsvHeader)
}

func (ds *DataSet) deriveContextFromRecords(inputRecords [][]string) *tableContext {
	context := tableContext{records: inputRecords}

	context.rowSize = uint(len(inputRecords)) - 1
	context.colSize = uint(len(inputRecords[0]))

	context.table = new(tables.CsvTableImpl)
	context.table.SetColumnAndRowSize(context.colSize, context.rowSize)

	return &context
}

func toBaseType(value string) interface{} {
	return caster.Cast(value)
}

func parseCsvText(csvText string) ([][]string, error) {
	textReader := strings.NewReader(csvText)

	r := csv.NewReader(textReader)
	r.TrimLeadingSpace = true

	records, readError := r.ReadAll()
	if readError != nil {
		return nil, errors.Wrap(readError, "parsing csv text")
	}

	return records, nil
}

func loadCsvRecords(filePath string) ([][]string, error) {
	fileHandle, openError := os.Open(filePath)
	defer fileHandle.Close()
	if openError != nil {
		return nil, errors.Wrap(openError, "opening csv file")
	}

	return loadCsvRecordsFromFileHandle(fileHandle)
}

func loadCsvRecordsFromFileHandle(fileHandle *os.File) ([][]string, error) {
	r := csv.NewReader(fileHandle)
	r.TrimLeadingSpace = true

	records, readError := r.ReadAll()
	if readError != nil {
		return nil, errors.Wrap(readError, "reading csv file")
	}

	return records, nil
}

func deriveHeader(records [][]string) dataset.TableHeader {
	newCsvHeader := make(dataset.TableHeader, 0)
	for _, heading := range records[0] {
		newCsvHeader = append(newCsvHeader, heading)
	}
	return newCsvHeader
}

func (ds *DataSet) deriveMetaTableName(metaCsvFilePath string) string {
	basePath := filepath.Base(metaCsvFilePath)

	splitPath := strings.Split(basePath, ".")
	extension := "." + splitPath[len(splitPath)-1]

	basePathSansExtension := strings.TrimSuffix(basePath, extension)
	return basePathSansExtension
}
