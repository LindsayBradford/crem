// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
)

type Cell struct {
	dispatch *ole.IDispatch
}

func (this *Cell) Value() interface{} {
	return this.getPropertyVariant("Value")
}

func (this *Cell) SetValue(value interface{}) {
	this.setProperty("Value", value)
}

func (this *Cell) getPropertyVariant(propertyName string, parameters ...interface{}) interface{} {
	return getPropertyVariant(this.dispatch, propertyName, parameters...)
}

func (this *Cell) setProperty(propertyName string, propertyValue interface{}) {
	setProperty(this.dispatch, propertyName, propertyValue)
}
