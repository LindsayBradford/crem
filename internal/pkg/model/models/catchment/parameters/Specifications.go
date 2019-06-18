// Copyright (c) 2019 Australian Rivers Institute.

package parameters

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters/specification"
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
	GullyRestorationCostPerKilometer     string = "GullyRestorationCostPerKilometer"
)

func DefineSpecifications() *specification.Specifications {
	specs := specification.NewSpecifications()
	specs.Add(
		specification.Specification{
			Key:          BankErosionFudgeFactor,
			Validator:    validateIsBankErosionFudgeFactor,
			DefaultValue: 5 * math.Pow(10, -4),
		},
	).Add(
		specification.Specification{
			Key:          WaterDensity,
			Validator:    specification.IsDecimal,
			DefaultValue: 1.0,
		},
	).Add(
		specification.Specification{
			Key:          LocalAcceleration,
			Validator:    specification.IsDecimal,
			DefaultValue: 9.81,
		},
	).Add(
		specification.Specification{
			Key:          GullyCompensationFactor,
			Validator:    specification.IsDecimal,
			DefaultValue: 0.5,
		},
	).Add(
		specification.Specification{
			Key:          SedimentDensity,
			Validator:    specification.IsDecimal,
			DefaultValue: 1.5,
		},
	).Add(
		specification.Specification{
			Key:          SuspendedSedimentProportion,
			Validator:    specification.IsDecimal,
			DefaultValue: 0.5,
		},
	).Add(
		specification.Specification{
			Key:          YearsOfErosion,
			Validator:    specification.IsNonNegativeInteger,
			DefaultValue: int64(100),
		},
	).Add(
		specification.Specification{
			Key:          DataSourcePath,
			Validator:    specification.IsReadableFile,
			DefaultValue: "",
		},
	).Add(
		specification.Specification{
			Key:          RiparianRevegetationCostPerKilometer,
			Validator:    specification.IsDecimal,
			DefaultValue: float64(24000),
		},
	).Add(
		specification.Specification{
			Key:          GullyRestorationCostPerKilometer,
			Validator:    specification.IsDecimal,
			DefaultValue: float64(44000),
		},
	)
	return specs
}

func validateIsBankErosionFudgeFactor(key string, value interface{}) error {
	minValue := 1 * math.Pow(10, -4)
	maxValue := 5 * math.Pow(10, -4)
	return specification.IsDecimalWithInclusiveBounds(key, value, minValue, maxValue)
}
