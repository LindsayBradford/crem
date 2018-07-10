// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func getProperty(iDispatch *ole.IDispatch, propertyName string, parameters... interface{}) *ole.IDispatch {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).ToIDispatch()
}

func getPropertyValue(iDispatch *ole.IDispatch, propertyName string, parameters... interface{}) int64 {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).Val
}

func getPropertyVariant(iDispatch *ole.IDispatch, propertyName string, parameters... interface{}) interface{} {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).Value()
}

func getPropertyString(iDispatch *ole.IDispatch, propertyName string, parameters... interface{}) string {
	return oleutil.MustGetProperty(iDispatch, propertyName, parameters...).ToString()
}

func setProperty(iDispatch *ole.IDispatch, propertyName string, propertyValue interface{}) {
	oleutil.PutProperty(iDispatch, propertyName, propertyValue)
}

func callMethod(dispatch *ole.IDispatch, methodName string, parameters... interface{}) *ole.IDispatch {
	return oleutil.MustCallMethod(dispatch, methodName, parameters...).ToIDispatch()
}
