// Copyright (c) 2019 Australian Rivers Institute.

package averaged

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	. "github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
)

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Parameters.
		Initialise("Suppapitnarm Coolant Parameter Validation").
		Enforcing(ParameterSpecifications())
	return p
}

const (
	_                          = iota
	CoolingFactor       string = "coolingFactor"
	StartingTemperature string = "StartingTemperature"
)

func ParameterSpecifications() *Specifications {
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
