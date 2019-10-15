// Copyright (c) 2019 Australian Rivers Institute.

package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	errors2 "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/pkg/errors"
)

func NewDataSet(name string) *DataSet {
	dataSet := new(DataSet)
	dataSet.Initialise(name)
	return dataSet
}

var metaTableHeadings = [2]string{"TableName", "FilePath"}

type DataSet struct {
	dataset.DataSetImpl
	filePath string
	errors   *errors2.CompositeError
}

func (ds *DataSet) Load(baseCsvFilePath string) error {
	ds.filePath = baseCsvFilePath
	ds.errors = errors2.New("Csv File Load Errors")
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
	table := new(tables.CsvTableImpl)

	records, loadError := loadCvsRecords(csvFilePath)
	if loadError != nil {
		ds.errors.Add(loadError)
		return nil
	}

	rowSize := uint(len(records)) - 1
	colSize := uint(len(records[0]))
	table.SetColumnAndRowSize(colSize, rowSize)

	newCsvHeader := deriveHeader(records)
	table.SetHeader(newCsvHeader)

	for rowIndex := uint(1); rowIndex <= rowSize; rowIndex++ {
		for colIndex := uint(0); colIndex < colSize; colIndex++ {
			recordAsBaseType := toBaseType(records[rowIndex][colIndex])
			table.SetCell(colIndex, rowIndex-1, recordAsBaseType)
		}
	}

	return table
}

func toBaseType(value string) interface{} {
	valueAsInt, intError := strconv.ParseInt(value, 10, 64)
	if intError == nil {
		return valueAsInt
	}
	valueAsUInt, uintError := strconv.ParseUint(value, 10, 64)
	if uintError == nil {
		return valueAsUInt
	}
	valueAsFloat, floatError := strconv.ParseFloat(value, 64)
	if floatError == nil {
		return valueAsFloat
	}
	valueAsBool, boolError := strconv.ParseBool(value)
	if boolError == nil {
		return valueAsBool
	}
	return value
}

func loadCvsRecords(filePath string) ([][]string, error) {
	metaFile, openError := os.Open(filePath)
	defer metaFile.Close()
	if openError != nil {
		return nil, errors.Wrap(openError, "opening csv meta-file")
	}

	r := csv.NewReader(metaFile)
	r.TrimLeadingSpace = true

	records, readError := r.ReadAll()
	if readError != nil {
		return nil, errors.Wrap(readError, "reading csv meta-file")
	}

	return records, nil
}

func deriveHeader(records [][]string) tables.CsvHeader {
	newCsvHeader := make(tables.CsvHeader, 0)
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
