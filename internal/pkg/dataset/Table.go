// Copyright (c) 2018 Australian Rivers Institute.

package dataset

type Table interface {
	Name() string
	SetName(name string)

	Cell(xPos uint, yPos uint) interface{}
	SetCell(xPos uint, yPos uint, value interface{})

	SetSize(colNum uint, rowNum uint)
}
