// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
)

type Cell interface {
	Value() interface{}
	SetValue(value interface{})
	SetNumberFormat(value interface{})
	Release()
}

type CellImpl struct {
	oleWrapper
}

func (cell *CellImpl) WithDispatch(dispatch *ole.IDispatch) *CellImpl {
	cell.dispatch = dispatch
	return cell
}

func (cell *CellImpl) Value() interface{} {
	return cell.getPropertyVariant("Value")
}

func (cell *CellImpl) SetValue(value interface{}) {
	cell.setProperty("Value", value)
}

func (cell *CellImpl) SetNumberFormat(value interface{}) {
	cell.setProperty("NumberFormat", value)
}

func (cell *CellImpl) getPropertyVariant(propertyName string, parameters ...interface{}) interface{} {
	return getPropertyVariant(cell.dispatch, propertyName, parameters...)
}

func (cell *CellImpl) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(cell.dispatch, propertyName, propertyValue)
}
