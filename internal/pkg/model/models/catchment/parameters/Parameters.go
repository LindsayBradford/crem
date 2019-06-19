// Copyright (c) 2018 Australian Rivers Institute.

package parameters

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"math"

	. "github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
)

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Parameters.
		Initialise("Catchment Model Parameter Validation").
		Enforcing(ParameterSpecifications())
	return p
}

const (
	BankErosionFudgeFactor      string = "BankErosionFudgeFactor"
	WaterDensity                string = "WaterDensity"
	LocalAcceleration           string = "LocalAcceleration"
	GullyCompensationFactor     string = "GullyCompensationFactor"
	SedimentDensity             string = "SedimentDensity"
	SuspendedSedimentProportion string = "SuspendedSedimentProportion"
	YearsOfErosion              string = "YearsOfErosion"
	DataSourcePath              string = "DataSourcePath"

	RiparianRevegetationCostPerKilometer string = "RiparianRevegetationCostPerKilometer"
	GullyRestorationCostPerKilometer     string = "GullyRestorationCostPerKilometer"
)

func ParameterSpecifications() *Specifications {
	specs := NewSpecifications()
	specs.Add(
		Specification{
			Key:          BankErosionFudgeFactor,
			Validator:    validateIsBankErosionFudgeFactor,
			DefaultValue: 5 * math.Pow(10, -4),
		},
	).Add(
		Specification{
			Key:          WaterDensity,
			Validator:    IsDecimal,
			DefaultValue: 1.0,
		},
	).Add(
		Specification{
			Key:          LocalAcceleration,
			Validator:    IsDecimal,
			DefaultValue: 9.81,
		},
	).Add(
		Specification{
			Key:          GullyCompensationFactor,
			Validator:    IsDecimal,
			DefaultValue: 0.5,
		},
	).Add(
		Specification{
			Key:          SedimentDensity,
			Validator:    IsDecimal,
			DefaultValue: 1.5,
		},
	).Add(
		Specification{
			Key:          SuspendedSedimentProportion,
			Validator:    IsDecimal,
			DefaultValue: 0.5,
		},
	).Add(
		Specification{
			Key:          YearsOfErosion,
			Validator:    IsNonNegativeInteger,
			DefaultValue: int64(100),
		},
	).Add(
		Specification{
			Key:          DataSourcePath,
			Validator:    IsReadableFile,
			DefaultValue: "",
		},
	).Add(
		Specification{
			Key:          RiparianRevegetationCostPerKilometer,
			Validator:    IsDecimal,
			DefaultValue: float64(24000),
		},
	).Add(
		Specification{
			Key:          GullyRestorationCostPerKilometer,
			Validator:    IsDecimal,
			DefaultValue: float64(44000),
		},
	)
	return specs
}

func validateIsBankErosionFudgeFactor(key string, value interface{}) error {
	minValue := 1 * math.Pow(10, -4)
	maxValue := 5 * math.Pow(10, -4)
	return IsDecimalWithInclusiveBounds(key, value, minValue, maxValue)
}
