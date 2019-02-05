// Copyright (c) 2018 Australian Rivers Institute.

package parameters

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
)

const (
	_                                  = iota
	BankErosionFudgeFactor      string = "BankErosionFudgeFactor"
	WaterDensity                string = "WaterDensity"
	LocalAcceleration           string = "LocalAcceleration"
	GullyCompensationFactor     string = "GullyCompensationFactor"
	SedimentDensity             string = "SedimentDensity"
	SuspendedSedimentProportion string = "SuspendedSedimentProportion"
	YearsOfErosion              string = "YearsOfErosion"
	DataSourcePath              string = "DataSourcePath"

	RiparianRevegetationCostPerKilometer string = "RiparianRevegetationCostPerKilometer"
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

	p.AddMetaData(
		parameters.MetaData{
			Key:          YearsOfErosion,
			Validator:    p.IsNonNegativeInteger,
			DefaultValue: 100,
		},
	)

	p.AddMetaData(
		parameters.MetaData{
			Key:          DataSourcePath,
			Validator:    p.IsReadableFile,
			DefaultValue: "",
		},
	)

	p.AddMetaData(
		parameters.MetaData{
			Key:          RiparianRevegetationCostPerKilometer,
			Validator:    p.IsDecimal,
			DefaultValue: float64(24000),
		},
	)
}

func (p *Parameters) validateIsBankErosionFudgeFactor(key string, value interface{}) bool {
	minValue := 1 * math.Pow(10, -4)
	maxValue := 5 * math.Pow(10, -4)
	return p.IsDecimalWithInclusiveBounds(key, value, minValue, maxValue)
}
