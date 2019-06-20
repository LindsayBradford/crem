// Copyright (c) 2019 Australian Rivers Institute.

package annealers

import (
	. "github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
)

const MaximumIterations string = "MaximumIterations"

func DefineSpecifications() *Specifications {
	specs := NewSpecifications()
	specs.Add(
		Specification{
			Key:          MaximumIterations,
			Validator:    IsNonNegativeInteger,
			DefaultValue: int64(0),
		},
	)
	return specs
}
