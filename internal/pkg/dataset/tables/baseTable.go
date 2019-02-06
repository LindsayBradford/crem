// Copyright (c) 2018 Australian Rivers Institute.

package tables

type baseTable struct {
	name string

	cells [][]interface{}
}

func (bt *baseTable) Name() string {
	return bt.name
}

func (bt *baseTable) SetName(name string) {
	bt.name = name
}

func (bt *baseTable) SetSize(colNum uint, rowNum uint) {
	bt.cells = make([][]interface{}, colNum)
	for col := uint(0); col < colNum; col++ {
		bt.cells[col] = make([]interface{}, rowNum)
	}
}

func (bt *baseTable) Size() (colNum uint, rowNum uint) {
	colNum = uint(len(bt.cells))
	firstRow := bt.cells[0]
	rowNum = uint(len(firstRow))
	return
}

func (bt *baseTable) Cell(xPos uint, yPos uint) interface{} {
	return bt.cells[xPos][yPos]
}

func (bt *baseTable) CellFloat64(xPos uint, yPos uint) float64 {
	return bt.cells[xPos][yPos].(float64)
}

func (bt *baseTable) CellString(xPos uint, yPos uint) string {
	return bt.cells[xPos][yPos].(string)
}

func (bt *baseTable) SetCell(xPos uint, yPos uint, value interface{}) {
	bt.cells[xPos][yPos] = value
}
