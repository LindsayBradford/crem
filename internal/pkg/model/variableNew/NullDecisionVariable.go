// Copyright (c) 2019 Australian Rivers Institute.

package variableNew

const NullDecisionVariableName = "NullDecisionVariable"

var NullDecisionVariable = new(nullDecisionVariable)

type nullDecisionVariable struct{}

func (ndv *nullDecisionVariable) Name() string                                 { return NullDecisionVariableName }
func (ndv *nullDecisionVariable) SetName(name string)                          {}
func (ndv *nullDecisionVariable) Value() float64                               { return 0 }
func (ndv *nullDecisionVariable) SetValue(value float64)                       {}
func (ndv *nullDecisionVariable) UnitOfMeasure() UnitOfMeasure                 { return "" }
func (ndv *nullDecisionVariable) SetUnitOfMeasure(unitOfMeasure UnitOfMeasure) {}
