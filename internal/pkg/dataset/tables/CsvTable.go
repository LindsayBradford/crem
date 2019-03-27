// Copyright (c) 2018 Australian Rivers Institute.

package tables

import "github.com/LindsayBradford/crem/internal/pkg/dataset"

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
