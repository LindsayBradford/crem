// Copyright (c) 2018 Australian Rivers Institute.

package tables

var DefaultNullTable = new(NullTable)

type NullTable struct{}

func (nt *NullTable) Name() string        { return "Null" }
func (nt *NullTable) SetName(name string) {}

func (nt *NullTable) Cell(xPos uint, yPos uint) interface{}           { return "null" }
func (nt *NullTable) CellFloat64(xPos uint, yPos uint) float64        { return float64(0) }
func (nt *NullTable) CellInt64(xPos uint, yPos uint) int64            { return int64(0) }
func (nt *NullTable) SetCell(xPos uint, yPos uint, value interface{}) {}

func (nt *NullTable) SetSize(colNum uint, rowNum uint) {}
func (nt *NullTable) Size() (colNum uint, rowNum uint) { return 0, 0 }
