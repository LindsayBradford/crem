// Copyright (c) 2018 Australian Rivers Institute.

package tables

var DefaultNullTable = new(NullTable)

type NullTable struct{}

func (nt *NullTable) Name() string        { return "Null" }
func (nt *NullTable) SetName(name string) {}

func (nt *NullTable) Cell(xPos uint, yPos uint) interface{}           { return "null" }
func (nt *NullTable) SetCell(xPos uint, yPos uint, value interface{}) {}

func (nt *NullTable) SetSize(colNum uint, rowNum uint) {}
