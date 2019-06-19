// Copyright (c) 2019 Australian Rivers Institute.

package parameters

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	. "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters/specification"
)

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Enforces(DefineSpecifications())
	return p
}

const (
	InitialObjectiveOneValue   = "InitialObjectiveOneValue"
	InitialObjectiveTwoValue   = "InitialObjectiveTwoValue"
	InitialObjectiveThreeValue = "InitialObjectiveThreeValue"

	NumberOfPlanningUnits = "NumberOfPlanningUnits"
)

func DefineSpecifications() *Specifications {
	specs := NewSpecifications()
	specs.Add(
		Specification{
			Key:          InitialObjectiveOneValue,
			Validator:    IsDecimal,
			DefaultValue: float64(1000),
		},
	).Add(
		Specification{
			Key:          InitialObjectiveTwoValue,
			Validator:    IsDecimal,
			DefaultValue: float64(2000),
		},
	).Add(
		Specification{
			Key:          InitialObjectiveThreeValue,
			Validator:    IsDecimal,
			DefaultValue: float64(3000),
		},
	).Add(
		Specification{
			Key:          NumberOfPlanningUnits,
			Validator:    IsNonNegativeInteger,
			DefaultValue: int64(100),
		},
	)
	return specs
}
