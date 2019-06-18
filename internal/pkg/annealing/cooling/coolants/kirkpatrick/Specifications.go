// Copyright (c) 2019 Australian Rivers Institute.

package kirkpatrick

import (
	. "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters/specification"
)

func DefineSpecifications() *Specifications {
	specs := NewSpecifications()
	specs.Add(
		Specification{
			Key:          StartingTemperature,
			Validator:    IsNonNegativeDecimal,
			DefaultValue: float64(0),
		},
	).Add(
		Specification{
			Key:          CoolingFactor,
			Validator:    IsDecimalBetweenZeroAndOne,
			DefaultValue: float64(1.0),
		},
	)
	return specs
}
