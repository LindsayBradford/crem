// Copyright (c) 2019 Australian Rivers Institute.

package variable

type DecisionVariables map[string]DecisionVariable

type DecisionVariable interface {
	Name() string
	Value() float64
}
