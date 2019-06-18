// Copyright (c) 2019 Australian Rivers Institute.

package dumb

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters/specification"
)

const (
	InitialObjectiveValue = "InitialObjectiveValue"
	MinimumObjectiveValue = "MinimumObjectiveValue"
	MaximumObjectiveValue = "MaximumObjectiveValue"
)

func DefineSpecifications() *specification.Specifications {
	specs := specification.New()
	specs.Add(
		specification.Specification{
			Key:          InitialObjectiveValue,
			Validator:    specification.IsDecimal,
			DefaultValue: float64(1000),
		},
	).Add(
		specification.Specification{
			Key:          MinimumObjectiveValue,
			Validator:    specification.IsDecimal,
			DefaultValue: float64(0),
		},
	).Add(
		specification.Specification{
			Key:          MaximumObjectiveValue,
			Validator:    specification.IsDecimal,
			DefaultValue: float64(2000),
		},
	)
	return specs
}
