// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"math"

	"github.com/LindsayBradford/crem/annealing/explorer"
)

const (
	_                                  = iota
	BankErosionFudgeFactor      string = "BankErosionFudgeFactor"
	WaterDensity                string = "WaterDensity"
	LocalAcceleration           string = "LocalAcceleration"
	GullyCompensationFactor     string = "GullyCompensationFactor"
	SedimentDensity             string = "SedimentDensity"
	SuspendedSedimentProportion string = "SuspendedSedimentProportion"
)

type SedimentTransportSolutionExplorer struct {
	explorer.SingleObjectiveAnnealableExplorer

	parameters *SedimentTransportParameters
}

func (stse *SedimentTransportSolutionExplorer) WithName(name string) *SedimentTransportSolutionExplorer {
	stse.SingleObjectiveAnnealableExplorer.WithName(name)
	return stse
}

func (stse *SedimentTransportSolutionExplorer) WithParameters(params map[string]interface{}) *SedimentTransportSolutionExplorer {
	stse.parameters = new(SedimentTransportParameters).Initialise()
	stse.parameters.Merge(params)
	stse.parameters.Validate()
	return stse
}

func (stse *SedimentTransportSolutionExplorer) ParameterErrors() error {
	return stse.parameters.ValidationErrors()
}

type SedimentTransportParameters struct {
	Parameters
}

func (p *SedimentTransportParameters) Initialise() *SedimentTransportParameters {
	p.Parameters.Initialise()
	p.buildParameterMetaData()
	p.CreateDefaultParameters()
	return p
}

func (p *SedimentTransportParameters) buildParameterMetaData() {
	p.metaDataMap[BankErosionFudgeFactor] = ParameterMetaData{
		validator:    p.validateIsBankErosionFudgeFactor,
		defaultValue: 5 * math.Pow(10, -4),
	}

	p.metaDataMap[WaterDensity] = ParameterMetaData{
		validator:    p.ValidateIsDecimal,
		defaultValue: 1.0,
	}

	p.metaDataMap[LocalAcceleration] = ParameterMetaData{
		validator:    p.ValidateIsDecimal,
		defaultValue: 9.81,
	}

	p.metaDataMap[GullyCompensationFactor] = ParameterMetaData{
		validator:    p.ValidateIsDecimal,
		defaultValue: 0.5,
	}

	p.metaDataMap[SedimentDensity] = ParameterMetaData{
		validator:    p.ValidateIsDecimal,
		defaultValue: 1.5,
	}

	p.metaDataMap[SuspendedSedimentProportion] = ParameterMetaData{
		validator:    p.ValidateIsDecimal,
		defaultValue: 0.5,
	}
}

func (p *Parameters) validateIsBankErosionFudgeFactor(key string, value interface{}) {
	minValue := 1 * math.Pow(10, -4)
	maxValue := 5 * math.Pow(10, -4)
	p.ValidateDecimalWithInclusiveBounds(key, value, minValue, maxValue)
}
