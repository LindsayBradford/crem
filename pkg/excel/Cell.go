// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"
)

const textNumberFormat = "@"

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
	switch value.(type) {
	case fmt.Stringer:
		valueAsStringer := value.(fmt.Stringer)
		cell.setProperty("Value", valueAsStringer.String())
	case bool, int, uint, int32, uint32, uint64, int64, float32, float64, string, nil:
		cell.setProperty("Value", value)
	default:
		panic(errors.New("Attempt to set cell value value of unsupported type"))
	}
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
