// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
)

type Cell interface {
	Value() interface{}
	SetValue(value interface{})
}

type CellImpl struct {
	dispatch *ole.IDispatch
}

func (cell *CellImpl) Value() interface{} {
	return cell.getPropertyVariant("Value")
}

func (cell *CellImpl) SetValue(value interface{}) {
	cell.setProperty("Value", value)
}

func (cell *CellImpl) getPropertyVariant(propertyName string, parameters ...interface{}) interface{} {
	return getPropertyVariant(cell.dispatch, propertyName, parameters...)
}

func (cell *CellImpl) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(cell.dispatch, propertyName, propertyValue)
}

// TODO: can I find a way not to expose this?
func (cell *CellImpl) oleDispatch() *ole.IDispatch {
	return cell.dispatch
}
