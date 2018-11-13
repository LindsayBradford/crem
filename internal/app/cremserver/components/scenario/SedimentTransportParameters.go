// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"math"

	"github.com/LindsayBradford/crem/annealing/parameters"
)

type SedimentTransportParameters struct {
	parameters.Parameters
}

func (p *SedimentTransportParameters) Initialise() *SedimentTransportParameters {
	p.Parameters.Initialise()
	p.buildMetaData()
	p.CreateDefaults()
	return p
}

func (p *SedimentTransportParameters) buildMetaData() {
	p.AddMetaData(
		parameters.MetaData{
			Key:          BankErosionFudgeFactor,
			Validator:    p.validateIsBankErosionFudgeFactor,
			DefaultValue: 5 * math.Pow(10, -4),
		},
	)

	p.AddMetaData(
		parameters.MetaData{
			Key:          WaterDensity,
			Validator:    p.IsDecimal,
			DefaultValue: 1.0,
		},
	)

	p.AddMetaData(
		parameters.MetaData{
			Key:          LocalAcceleration,
			Validator:    p.IsDecimal,
			DefaultValue: 9.81,
		},
	)

	p.AddMetaData(
		parameters.MetaData{
			Key:          GullyCompensationFactor,
			Validator:    p.IsDecimal,
			DefaultValue: 0.5,
		},
	)

	p.AddMetaData(
		parameters.MetaData{
			Key:          SedimentDensity,
			Validator:    p.IsDecimal,
			DefaultValue: 1.5,
		},
	)

	p.AddMetaData(
		parameters.MetaData{
			Key:          SuspendedSedimentProportion,
			Validator:    p.IsDecimal,
			DefaultValue: 0.5,
		},
	)
}

func (p *SedimentTransportParameters) validateIsBankErosionFudgeFactor(key string, value interface{}) bool {
	minValue := 1 * math.Pow(10, -4)
	maxValue := 5 * math.Pow(10, -4)
	return p.IsDecimalWithInclusiveBounds(key, value, minValue, maxValue)
}
