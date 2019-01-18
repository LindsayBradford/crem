// Copyright (c) 2018 Australian Rivers Institute.

package tables

type CsvHeader []string

// var _ dataset.Table = (*CsvTable)(nil)
type CsvTable struct {
	baseTable
	header CsvHeader
}

func (ct *CsvTable) Header() CsvHeader {
	return ct.header
}

func (ct *CsvTable) SetHeader(header CsvHeader) {
	ct.header = header
}
