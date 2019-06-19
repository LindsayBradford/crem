// Copyright (c) 2018 Australian Rivers Institute.

package dumb

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	. "github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
)

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Enforces(ParameterSpecifications())
	return p
}

const (
	InitialObjectiveValue = "InitialObjectiveValue"
	MinimumObjectiveValue = "MinimumObjectiveValue"
	MaximumObjectiveValue = "MaximumObjectiveValue"
)

func ParameterSpecifications() *Specifications {
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
