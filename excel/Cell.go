// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
)

type Cell struct {
	dispatch *ole.IDispatch
}

func (cell *Cell) Value() interface{} {
	return cell.getPropertyVariant("Value")
}

func (cell *Cell) SetValue(value interface{}) {
	cell.setProperty("Value", value)
}

func (cell *Cell) getPropertyVariant(propertyName string, parameters ...interface{}) interface{} {
	return getPropertyVariant(cell.dispatch, propertyName, parameters...)
}

func (cell *Cell) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(cell.dispatch, propertyName, propertyValue)
}
