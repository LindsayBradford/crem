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

	RiparianBufferVegetationProportionTarget string = "RiparianBufferVegetationProportionTarget"
	GullySedimentReductionTarget             string = "GullySedimentReductionTarget"

	HillSlopeDeliveryRatio string = "HillSlopeDeliveryRatio"

	MaximumSedimentProduction            = "MaximumSedimentProduction"
	MaximumImplementationCost            = "MaximumImplementationCost"
	MaximumOpportunityCost               = "MaximumOpportunityCost"
	MaximumParticulateNitrogenProduction = "MaximumParticulateNitrogenProduction"
	MaximumDissolvedNitrogenProduction   = "MaximumDissolvedNitrogenProduction"
	MaximumTotalNitrogenProduction       = "MaximumTotalNitrogenProduction"
)

func ParameterSpecifications() *Specifications {
	specs := NewSpecifications()
	specs.Add(
		Specification{
			Key:          DataSourcePath,
			Validator:    IsReadableFile,
			DefaultValue: "",
		},
	).Add(
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
			Key:          RiparianBufferVegetationProportionTarget,
			Validator:    IsDecimalBetweenZeroAndOne,
			DefaultValue: float64(0.75),
		},
	).Add(
		Specification{
			Key:          GullySedimentReductionTarget,
			Validator:    IsDecimalBetweenZeroAndOne,
			DefaultValue: float64(0.8),
		},
	).Add(
		Specification{
			Key:          HillSlopeDeliveryRatio,
			Validator:    IsDecimalBetweenZeroAndOne,
			DefaultValue: float64(0.05),
		},
	).Add(
		Specification{
			Key:        MaximumSedimentProduction,
			Validator:  IsNonNegativeDecimal,
			IsOptional: true,
		},
	).Add(
		Specification{
			Key:        MaximumParticulateNitrogenProduction,
			Validator:  IsNonNegativeDecimal,
			IsOptional: true,
		},
	).Add(
		Specification{
			Key:        MaximumDissolvedNitrogenProduction,
			Validator:  IsNonNegativeDecimal,
			IsOptional: true,
		},
	).Add(
		Specification{
			Key:        MaximumTotalNitrogenProduction,
			Validator:  IsNonNegativeDecimal,
			IsOptional: true,
		},
	).Add(
		Specification{
			Key:        MaximumImplementationCost,
			Validator:  IsNonNegativeDecimal,
			IsOptional: true,
		},
	).Add(
		Specification{
			Key:        MaximumOpportunityCost,
			Validator:  IsNonNegativeDecimal,
			IsOptional: true,
		},
	)

	return specs
}

func validateIsBankErosionFudgeFactor(key string, value interface{}) error {
	minValue := 1 * math.Pow(10, -4)
	maxValue := 5 * math.Pow(10, -4)
	return IsDecimalWithInclusiveBounds(key, value, minValue, maxValue)
}
