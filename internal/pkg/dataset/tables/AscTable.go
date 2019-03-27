// Copyright (c) 2018 Australian Rivers Institute.

package tables

import "github.com/LindsayBradford/crem/internal/pkg/dataset"

var _ AscTable = new(AscTableImpl)

type AscTable interface {
	dataset.Table
	Header() AscHeader
	SetHeader(header AscHeader)
}

type AscHeader struct {
	NumCols     uint
	NumRows     uint
	XllCorner   float64
	YllCorner   float64
	CellSize    int64
	NoDataValue int64
}

type AscTableImpl struct {
	baseTable
	header AscHeader
}

func (at *AscTableImpl) Header() AscHeader {
	return at.header
}

func (at *AscTableImpl) SetHeader(header AscHeader) {
	at.header = header
}
