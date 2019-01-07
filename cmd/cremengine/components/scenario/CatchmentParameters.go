// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/dumb"
)

const (
	_                                  = iota
	BankErosionFudgeFactor      string = "BankErosionFudgeFactor"
	WaterDensity                string = "WaterDensity"
	LocalAcceleration           string = "LocalAcceleration"
	GullyCompensationFactor     string = "GullyCompensationFactor"
	SedimentDensity             string = "SedimentDensity"
	SuspendedSedimentProportion string = "SuspendedSedimentProportion"
	DataSourcePath              string = "DataSourcePath"
)

type CatchmentParameters struct {
	parameters.Parameters
}

func (cp *CatchmentParameters) Initialise() *CatchmentParameters {
	cp.Parameters.Initialise()
	cp.buildMetaData()
	cp.CreateDefaults()
	return cp
}

func (cp *CatchmentParameters) buildMetaData() {
	cp.AddMetaData( // TODO: This is just a placeholder for initial testing
		parameters.MetaData{
			Key:          dumb.InitialObjectiveValue,
			Validator:    cp.IsNonNegativeDecimal,
			DefaultValue: 2000.0,
		},
	)

	cp.AddMetaData(
		parameters.MetaData{
			Key:          BankErosionFudgeFactor,
			Validator:    cp.validateIsBankErosionFudgeFactor,
			DefaultValue: 5 * math.Pow(10, -4),
		},
	)

	cp.AddMetaData(
		parameters.MetaData{
			Key:          WaterDensity,
			Validator:    cp.IsDecimal,
			DefaultValue: 1.0,
		},
	)

	cp.AddMetaData(
		parameters.MetaData{
			Key:          LocalAcceleration,
			Validator:    cp.IsDecimal,
			DefaultValue: 9.81,
		},
	)

	cp.AddMetaData(
		parameters.MetaData{
			Key:          GullyCompensationFactor,
			Validator:    cp.IsDecimal,
			DefaultValue: 0.5,
		},
	)

	cp.AddMetaData(
		parameters.MetaData{
			Key:          SedimentDensity,
			Validator:    cp.IsDecimal,
			DefaultValue: 1.5,
		},
	)

	cp.AddMetaData(
		parameters.MetaData{
			Key:          SuspendedSedimentProportion,
			Validator:    cp.IsDecimal,
			DefaultValue: 0.5,
		},
	)

	cp.AddMetaData(
		parameters.MetaData{
			Key:          DataSourcePath,
			Validator:    cp.IsReadableFile,
			DefaultValue: "",
		},
	)
}

func (cp *CatchmentParameters) validateIsBankErosionFudgeFactor(key string, value interface{}) bool {
	minValue := 1 * math.Pow(10, -4)
	maxValue := 5 * math.Pow(10, -4)
	return cp.IsDecimalWithInclusiveBounds(key, value, minValue, maxValue)
}