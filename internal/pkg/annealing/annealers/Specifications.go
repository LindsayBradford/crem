// Copyright (c) 2019 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters/specification"
)

const MaximumIterations string = "MaximumIterations"

func DefineSpecifications() *specification.Specifications {
	specs := specification.NewSpecifications()
	specs.Add(
		specification.Specification{
			Key:          MaximumIterations,
			Validator:    specification.IsNonNegativeInteger,
			DefaultValue: int64(0),
		},
	)
	return specs
}
