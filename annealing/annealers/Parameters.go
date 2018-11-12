// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import "github.com/LindsayBradford/crem/annealing/parameters"

const (
	_                          = iota
	CoolingFactor       string = "CoolingFactor"
	MaximumIterations   string = "MaximumIterations"
	StartingTemperature string = "StartingTemperature"
)

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Parameters.Initialise()
	p.buildMetaData()
	p.CreateDefaults()
	return p
}

func (p *Parameters) buildMetaData() {
	p.AddMetaData(
		parameters.MetaData{
			Key:          StartingTemperature,
			Validator:    p.ValidateIsNonNegativeDecimal,
			DefaultValue: float64(0),
		},
	)

	p.AddMetaData(
		parameters.MetaData{
			Key:          MaximumIterations,
			Validator:    p.ValidateIsUnsignedInteger,
			DefaultValue: int64(0),
		},
	)

	p.AddMetaData(
		parameters.MetaData{
			Key:          CoolingFactor,
			Validator:    p.ValidateIsDecimalBetweenZeroAndOne,
			DefaultValue: float64(1.0),
		},
	)
}
