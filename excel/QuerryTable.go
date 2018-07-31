// Copyright (c) 2018 Australian Rivers Institute.

package excel

import "github.com/go-ole/go-ole"

type QueryTable interface {
	SetProperty(propertyName string, propertyValue interface{})
}

type QueryTableImpl struct {
	dispatch *ole.IDispatch
}

func (qt *QueryTableImpl) SetProperty(propertyName string, propertyValue interface{}) {
	setProperty(qt.dispatch, propertyName, propertyValue)
}
