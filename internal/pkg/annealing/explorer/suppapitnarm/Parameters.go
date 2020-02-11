// Copyright (c) 2018 Australian Rivers Institute.

package suppapitnarm

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters"

	. "github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
)

const DefaultExplorableDecisionVariables = "SedimentProduced,ImplementationCost"

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Parameters.
		Initialise("Suppapitnarm Explorer Parameter Validation").
		Enforcing(ParameterSpecifications())
	return p
}

const (
	ExplorableDecisionVariables   = "ExplorableDecisionVariables"
	ReturnToBaseAdjustmentFactor  = "ReturnToBaseAdjustmentFactor"
	InitialReturnToBaseStep       = "InitialReturnToBaseStep"
	MinimumReturnToBaseRate       = "MinimumReturnToBaseRate"
	ReturnToBaseIsolationFraction = "ReturnToBaseIsolationFraction"
)

type optimisationDirection int

const (
	Invalid optimisationDirection = iota
	Minimising
	Maximising
)

func ParameterSpecifications() *Specifications {
	specs := NewSpecifications()
	specs.Add(
		Specification{
			Key:          ExplorableDecisionVariables,
			Validator:    IsString,
			DefaultValue: DefaultExplorableDecisionVariables,
		},
	).Add(
		Specification{
			Key:          ReturnToBaseAdjustmentFactor,
			Validator:    IsDecimalBetweenZeroAndOne,
			DefaultValue: float64(0.95),
		},
	).Add(
		Specification{
			Key:          InitialReturnToBaseStep,
			Validator:    IsNonNegativeInteger,
			DefaultValue: int64(20_000), // following initial CRP hard-coded default
		},
	).Add(
		Specification{
			Key:          MinimumReturnToBaseRate,
			Validator:    IsNonNegativeInteger,
			DefaultValue: int64(10), // following initial CRP hard-coded default
		},
	).Add(
		Specification{
			Key:          ReturnToBaseIsolationFraction,
			Validator:    IsDecimalBetweenZeroAndOne,
			DefaultValue: float64(0.9), // following initial CRP hard-coded default
		},
	)
	return specs
}
