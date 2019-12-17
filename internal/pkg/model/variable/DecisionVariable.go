// Copyright (c) 2019 Australian Rivers Institute.

// Variable package supplies generalised model variables that allow watchers of models to make decisions on those
// models by reacting the changes in model decision variables. Decision variables for a model are considered part of
// that model's public interface.
package variable

import (
	"sort"

	"github.com/LindsayBradford/crem/pkg/name"
)

type DecisionVariableMap map[string]DecisionVariable

func (m DecisionVariableMap) SortedKeys() (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

// DecisionVariable describes an interface between a Model and any decision making logic observing the model via its
// decision variables.  This Value of a decision Variable should be a fine-grained indicator of how well a model is
// doing against some objective we have for that model.
// There should be one decision Variable representing each objective being evaluated for a model.
type DecisionVariable interface {
	// Name returns the model-centric name of the DecisionVariable.
	// Decision variables are expected to have unique names within a model.
	name.Nameable

	Value() float64
	SetValue(value float64)

	UnitOfMeasure() UnitOfMeasure
	Precision() Precision
}

type UnitOfMeasure string

func (uom UnitOfMeasure) String() string { return string(uom) }

const (
	NotApplicable UnitOfMeasure = "Not Applicable (NA)"
	TonnesPerYear UnitOfMeasure = "Tonnes per Year (t/y)"
	Dollars       UnitOfMeasure = "Dollars ($)"
)

const defaultUnitOfMeasure = NotApplicable

type Precision int

const defaultPrecision = 3

func NewSimpleDecisionVariable(name string) SimpleDecisionVariable {
	return SimpleDecisionVariable{
		name:          name,
		value:         0,
		unitOfMeasure: defaultUnitOfMeasure,
		precision:     defaultPrecision,
	}
}

type SimpleDecisionVariable struct {
	name  string
	value float64

	unitOfMeasure UnitOfMeasure
	precision     Precision
}

func (v *SimpleDecisionVariable) Name() string                           { return v.name }
func (v *SimpleDecisionVariable) SetName(name string)                    { v.name = name }
func (v *SimpleDecisionVariable) Value() float64                         { return v.value }
func (v *SimpleDecisionVariable) SetValue(value float64)                 { v.value = value }
func (v *SimpleDecisionVariable) UnitOfMeasure() UnitOfMeasure           { return v.unitOfMeasure }
func (v *SimpleDecisionVariable) SetUnitOfMeasure(measure UnitOfMeasure) { v.unitOfMeasure = measure }
func (v *SimpleDecisionVariable) Precision() Precision                   { return v.precision }
func (v *SimpleDecisionVariable) SetPrecision(precision Precision)       { v.precision = precision }
