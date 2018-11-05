// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/errors"
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
	stse.parameters = params
	stse.parametersMetaData = new(ParametersMetaData).Initialise()
	stse.validateParameters()
	return stse
}

func (stse *SedimentTransportSolutionExplorer) validateParameters() {
	for key, value := range stse.parameters {
		stse.parametersMetaData.Validate(key, value)
	}
}

func (stse *SedimentTransportSolutionExplorer) ParameterErrors() error {
	return stse.parametersMetaData.Errors()
}

type ParametersMetaData struct {
	parameterMap map[string]parameterValidator
	errors       *errors.CompositeError
}

func (pmd *ParametersMetaData) Initialise() *ParametersMetaData {
	pmd.errors = errors.NewComposite("SedimentTransportSolutionExplorer parameters")
	pmd.buildParameterMetaData()
	return pmd
}

func (pmd *ParametersMetaData) buildParameterMetaData() {
	pmd.parameterMap = make(map[string]parameterValidator, 0)

	pmd.parameterMap["Penalty"] = pmd.validateValueIsDecimal

}

func (pmd *ParametersMetaData) validateValueIsDecimal(key string, value interface{}) {
	_, typeIsOk := value.(float64)
	if !typeIsOk {
		pmd.errors.AddMessage("Parameter [" + key + "] must be a decimal value")
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
		pmd.parameterMap[key](key, value)
	} else {
		pmd.keyMissingValidator(key)
	}
}

func (pmd *ParametersMetaData) keyMissingValidator(key string) {
	pmd.errors.AddMessage("Parameter [" + key + "] is not a parameter for this explorer")
}

type parameterValidator func(key string, value interface{})
