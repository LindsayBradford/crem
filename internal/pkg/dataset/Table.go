// Copyright (c) 2018 Australian Rivers Institute.

package dataset

type TableType int

type Table interface {
	Name() string
	SetName(name string)

	Cell(xPos uint, yPos uint) interface{}
	CellFloat64(xPos uint, yPos uint) float64
	CellString(xPos uint, yPos uint) string
	SetCell(xPos uint, yPos uint, value interface{})

	SetColumnAndRowSize(colNum uint, rowNum uint)
	ColumnAndRowSize() (colNum uint, rowNum uint)
}
