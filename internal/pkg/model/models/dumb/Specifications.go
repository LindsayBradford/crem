// Copyright (c) 2019 Australian Rivers Institute.

package dumb

import (
	. "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters/specification"
)

const (
	InitialObjectiveValue = "InitialObjectiveValue"
	MinimumObjectiveValue = "MinimumObjectiveValue"
	MaximumObjectiveValue = "MaximumObjectiveValue"
)

func DefineSpecifications() *Specifications {
	specs := NewSpecifications()
	specs.Add(
		Specification{
			Key:          InitialObjectiveValue,
			Validator:    IsDecimal,
			DefaultValue: float64(1000),
		},
	).Add(
		Specification{
			Key:          MinimumObjectiveValue,
			Validator:    IsDecimal,
			DefaultValue: float64(0),
		},
	).Add(
		Specification{
			Key:          MaximumObjectiveValue,
			Validator:    IsDecimal,
			DefaultValue: float64(2000),
		},
	)
	return specs
}
