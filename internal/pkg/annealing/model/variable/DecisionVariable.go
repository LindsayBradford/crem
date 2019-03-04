// Copyright (c) 2019 Australian Rivers Institute.

// variable package supplies generalised model variables that allow watchers of models to make decisions based on the
// state of those models by reacting the changes in their variables.
package variable

type DecisionVariables map[string]DecisionVariable

// DecisionVariable describes an interface between a Model and any decision makers observing the model that may
// alter their decision based on the value reported by the decision variable.
type DecisionVariable interface {
	Name() string
	Value() float64
}
