// Copyright (c) 2018 Australian Rivers Institute.

package tables

type AscHeader struct {
	NumCols     uint
	NumRows     uint
	XllCorner   float64
	YllCorner   float64
	CellSize    int64
	NoDataValue int64
}

type AscTable struct {
	baseTable
	header AscHeader
}

func (at *AscTable) Header() AscHeader {
	return at.header
}

func (at *AscTable) SetHeader(header AscHeader) {
	at.header = header
}
