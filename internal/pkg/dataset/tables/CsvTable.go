// Copyright (c) 2018 Australian Rivers Institute.

package tables

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/pkg/errors"
)

var _ CsvTable = new(CsvTableImpl)

type CsvTable interface {
	dataset.Table
	Header() CsvHeader
	SetHeader(header CsvHeader)
}

type CsvHeader []string

type CsvTableImpl struct {
	baseTable
	header CsvHeader
}

func (ct *CsvTableImpl) Header() CsvHeader {
	return ct.header
}

func (ct *CsvTableImpl) SetHeader(header CsvHeader) {
	ct.header = header
}

func ToCsvTable(dataSet dataset.DataSet, tableName string) CsvTable {
	namedTable, namedTableError := dataSet.Table(tableName)
	if namedTableError != nil {
		panic(errors.New("Expected data set supplied to have a [" + tableName + "] table"))
	}

	namedCsvTable, isCsvType := namedTable.(CsvTable)
	if !isCsvType {
		panic(errors.New("Expected data set table [" + tableName + "] to be a CSV type"))
	}
	return namedCsvTable
}
