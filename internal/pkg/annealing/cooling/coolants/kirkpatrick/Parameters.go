// Copyright (c) 2018 Australian Rivers Institute.

package kirkpatrick

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	. "github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
)

type Parameters struct {
	parameters.Parameters
}

func (kp *Parameters) Initialise() *Parameters {
	kp.Enforces(ParameterSpecifications())
	return kp
}

const (
	_                          = iota
	CoolingFactor       string = "CoolingFactor"
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
