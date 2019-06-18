// Copyright (c) 2019 Australian Rivers Institute.

package parameters

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters/specification"
)

const (
	InitialObjectiveOneValue   = "InitialObjectiveOneValue"
	InitialObjectiveTwoValue   = "InitialObjectiveTwoValue"
	InitialObjectiveThreeValue = "InitialObjectiveThreeValue"

	NumberOfPlanningUnits = "NumberOfPlanningUnits"
)

func DefineSpecifications() *specification.Specifications {
	specs := specification.New()
	specs.Add(
		specification.Specification{
			Key:          InitialObjectiveOneValue,
			Validator:    specification.IsDecimal,
			DefaultValue: float64(1000),
		},
	).Add(
		specification.Specification{
			Key:          InitialObjectiveTwoValue,
			Validator:    specification.IsDecimal,
			DefaultValue: float64(2000),
		},
	).Add(
		specification.Specification{
			Key:          InitialObjectiveThreeValue,
			Validator:    specification.IsDecimal,
			DefaultValue: float64(3000),
		},
	).Add(
		specification.Specification{
			Key:          NumberOfPlanningUnits,
			Validator:    specification.IsNonNegativeInteger,
			DefaultValue: int64(100),
		},
	)
	return specs
}
