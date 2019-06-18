// Copyright (c) 2019 Australian Rivers Institute.

package kirkpatrick

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters/specification"
)

func DefineSpecifications() *specification.Specifications {
	specs := specification.New()
	specs.Add(
		specification.Specification{
			Key:          StartingTemperature,
			Validator:    specification.IsNonNegativeDecimal,
			DefaultValue: float64(0),
		},
	).Add(
		specification.Specification{
			Key:          CoolingFactor,
			Validator:    specification.IsDecimalBetweenZeroAndOne,
			DefaultValue: float64(1.0),
		},
	)
	return specs
}
