// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"fmt"
	"math"

	"github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/errors"
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

	parameters         map[string]interface{}
	parametersMetaData *ParametersMetaData
}

func (stse *SedimentTransportSolutionExplorer) WithName(name string) *SedimentTransportSolutionExplorer {
	stse.SingleObjectiveAnnealableExplorer.WithName(name)
	return stse
}

func (stse *SedimentTransportSolutionExplorer) WithParameters(params map[string]interface{}) *SedimentTransportSolutionExplorer {
	stse.createDefaultParamsFromMetadata()
	stse.mergeParamsWithDefaults(params)
	stse.validateParams()
	return stse
}

func (stse *SedimentTransportSolutionExplorer) mergeParamsWithDefaults(params map[string]interface{}) {
	for suppliedKey, suppliedValue := range params {
		stse.parameters[suppliedKey] = suppliedValue
	}
}

func (stse *SedimentTransportSolutionExplorer) createDefaultParamsFromMetadata() {
	stse.parametersMetaData = new(ParametersMetaData).Initialise()
	stse.parameters = make(map[string]interface{}, 0)
	for key, value := range stse.parametersMetaData.parameterMap {
		stse.parameters[key] = value.defaultValue
	}
}

func (stse *SedimentTransportSolutionExplorer) validateParams() {
	for key, value := range stse.parameters {
		stse.parametersMetaData.Validate(key, value)
	}
}

func (stse *SedimentTransportSolutionExplorer) ParameterErrors() error {
	return stse.parametersMetaData.Errors()
}

type ParametersMetaData struct {
	parameterMap map[string]ParameterMetaData
	errors       *errors.CompositeError
}

type ParameterMetaData struct {
	validator    parameterValidator
	defaultValue interface{}
}

func (pmd *ParametersMetaData) Initialise() *ParametersMetaData {
	pmd.errors = errors.NewComposite("SedimentTransportSolutionExplorer parameters")
	pmd.buildParameterMetaData()
	return pmd
}

func (pmd *ParametersMetaData) buildParameterMetaData() {
	pmd.parameterMap = make(map[string]ParameterMetaData, 0)

	pmd.parameterMap[BankErosionFudgeFactor] = ParameterMetaData{
		validator:    pmd.validateIsBankErosionFudgeFactor,
		defaultValue: 5 * math.Pow(10, -4),
	}

	pmd.parameterMap[WaterDensity] = ParameterMetaData{
		validator:    pmd.validateIsDecimal,
		defaultValue: 1.0,
	}

	pmd.parameterMap[LocalAcceleration] = ParameterMetaData{
		validator:    pmd.validateIsDecimal,
		defaultValue: 9.81,
	}

	pmd.parameterMap[GullyCompensationFactor] = ParameterMetaData{
		validator:    pmd.validateIsDecimal,
		defaultValue: 0.5,
	}

	pmd.parameterMap[SedimentDensity] = ParameterMetaData{
		validator:    pmd.validateIsDecimal,
		defaultValue: 1.5,
	}

	pmd.parameterMap[SuspendedSedimentProportion] = ParameterMetaData{
		validator:    pmd.validateIsDecimal,
		defaultValue: 0.5,
	}
}

func (pmd *ParametersMetaData) validateIsDecimal(key string, value interface{}) {
	_, typeIsOk := value.(float64)
	if !typeIsOk {
		pmd.errors.AddMessage("Parameter [" + key + "] must be a decimal value")
	}
}

func (pmd *ParametersMetaData) validateIsBankErosionFudgeFactor(key string, value interface{}) {
	minValue := 1 * math.Pow(10, -4)
	maxValue := 5 * math.Pow(10, -4)
	pmd.validateDecimalWithInclusiveBounds(key, value, minValue, maxValue)
}

func (pmd *ParametersMetaData) validateDecimalWithInclusiveBounds(key string, value interface{}, minValue float64, maxValue float64) {
	valueAsFloat, typeIsOk := value.(float64)
	if !typeIsOk {
		pmd.errors.AddMessage("Parameter [" + key + "] must be a decimal value")
		return
	}

	if valueAsFloat < minValue || valueAsFloat > maxValue {
		message := fmt.Sprintf("Parameter [%s] supplied with decimal value [%v], but must be between [%.04f] and [%.04f] inclusive", key, value, minValue, maxValue)
		pmd.errors.AddMessage(message)
	}
}

func (pmd *ParametersMetaData) Errors() error {
	if pmd.errors.Size() > 0 {
		return pmd.errors
	}
	return nil
}

func (pmd *ParametersMetaData) Validate(key string, value interface{}) {
	if _, isPresent := pmd.parameterMap[key]; isPresent {
		pmd.parameterMap[key].validator(key, value)
	} else {
		pmd.keyMissingValidator(key)
	}
}

func (pmd *ParametersMetaData) keyMissingValidator(key string) {
	pmd.errors.AddMessage("Parameter [" + string(key) + "] is not a parameter for this explorer")
}

type parameterValidator func(key string, value interface{})
