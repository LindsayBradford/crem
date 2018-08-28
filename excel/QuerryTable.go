// Copyright (c) 2018 Australian Rivers Institute.

package excel

import "github.com/go-ole/go-ole"

type QueryTable interface {
	SetProperty(propertyName string, propertyValue interface{})
	Release()
}

type QueryTableImpl struct {
	oleWrapper
}

func (qt *QueryTableImpl) WithDispatch(dispatch *ole.IDispatch) *QueryTableImpl {
	qt.dispatch = dispatch
	return qt
}

func (qt *QueryTableImpl) SetProperty(propertyName string, propertyValue interface{}) {
	setProperty(qt.dispatch, propertyName, propertyValue)
}
