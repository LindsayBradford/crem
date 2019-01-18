// Copyright (c) 2018 Australian Rivers Institute.

package dataset

type Table interface {
	Name() string
	SetName(name string)

	Cell(xPos uint, yPos uint) interface{}
	CellFloat64(xPos uint, yPos uint) float64
	CellInt64(xPos uint, yPos uint) int64
	SetCell(xPos uint, yPos uint, value interface{})

	SetSize(colNum uint, rowNum uint)
	Size() (colNum uint, rowNum uint)
}
