// Copyright (c) 2019 Australian Rivers Institute.

// variable package supplies generalised model variables that allow watchers of models to make decisions on those
// models by reacting the changes in model decision variables. Decision variables for a model are considered part of
// that model's public interface.
package variable

import "github.com/LindsayBradford/crem/pkg/name"

type DecisionVariables map[string]DecisionVariable

// DecisionVariable describes an interface between a Model and any decision making logic observing the model via its
// decision variables.  This Value of a decision variable should be a fine-grained indicator of how well a model is
// doing against some objective we have for that model.
// There should be one decision variable representing each objective being evaluated for a model.
type DecisionVariable interface {
	// Name returns the model-centric name of the DecisionVariable.
	// Decision variables are expected to have unique names within a model.
	name.Nameable

	Value() float64
	SetValue(value float64)

	UnitOfMeasureContainer
}
