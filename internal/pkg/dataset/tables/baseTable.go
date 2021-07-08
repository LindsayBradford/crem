// Copyright (c) 2018 Australian Rivers Institute.

package tables

import (
	"fmt"
)

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

func (bt *baseTable) SetColumnAndRowSize(colNum uint, rowNum uint) {
	bt.cells = make([][]interface{}, rowNum)
	for row := uint(0); row < rowNum; row++ {
		bt.cells[row] = make([]interface{}, colNum)
	}
}

func (bt *baseTable) ColumnAndRowSize() (colNum uint, rowNum uint) {
	rowNum = uint(len(bt.cells))
	firstRow := bt.cells[0]
	colNum = uint(len(firstRow))
	return
}

func (bt *baseTable) Cell(col uint, row uint) interface{} {
	return bt.cells[row][col]
}

func (bt *baseTable) CellFloat64(col uint, row uint) float64 {
	return bt.cells[row][col].(float64)
}

func (bt *baseTable) CellString(col uint, row uint) string {
	switch bt.cells[row][col].(type) {
	case string:
		return bt.cells[row][col].(string)
	case float64:
		asFloat64 := bt.CellFloat64(col, row)
		return fmt.Sprintf("%v", asFloat64)
	default:
		return ""
	}
}

func (bt *baseTable) SetCell(col uint, row uint, value interface{}) {
	bt.cells[row][col] = value
}
