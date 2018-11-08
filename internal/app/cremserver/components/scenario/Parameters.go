// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"fmt"

	"github.com/LindsayBradford/crem/errors"
)

type ParameterMap map[string]interface{}
type ParameterMetaDataMap map[string]ParameterMetaData

type Parameters struct {
	parameterMap     ParameterMap
	metaDataMap      ParameterMetaDataMap
	validationErrors *errors.CompositeError
}

type ParameterValidator func(key string, value interface{})

type ParameterMetaData struct {
	validator    ParameterValidator
	defaultValue interface{}
}

func (p *Parameters) Initialise() *Parameters {
	p.validationErrors = errors.NewComposite("SolutionExplorer Parameters")
	p.metaDataMap = make(ParameterMetaDataMap, 0)
	return p
}

func (p *Parameters) CreateDefaultParameters() {
	p.parameterMap = make(ParameterMap, 0)
	for key, value := range p.metaDataMap {
		p.parameterMap[key] = value.defaultValue
	}
}

func (p *Parameters) Merge(params ParameterMap) {
	for suppliedKey, suppliedValue := range params {
		p.parameterMap[suppliedKey] = suppliedValue
	}
}

func (p *Parameters) ValidationErrors() error {
	if p.validationErrors.Size() > 0 {
		return p.validationErrors
	}
	return nil
}

func (p *Parameters) Validate() {
	for key, value := range p.parameterMap {
		p.validateParam(key, value)
	}
}

func (p *Parameters) validateParam(key string, value interface{}) {
	if _, isPresent := p.metaDataMap[key]; isPresent {
		p.metaDataMap[key].validator(key, value)
	} else {
		p.keyMissingValidator(key)
	}
}

func (p *Parameters) keyMissingValidator(key string) {
	p.validationErrors.AddMessage("Parameter [" + string(key) + "] is not a parameter for this explorer")
}

func (p *Parameters) ValidateIsDecimal(key string, value interface{}) {
	_, typeIsOk := value.(float64)
	if !typeIsOk {
		p.validationErrors.AddMessage("Parameter [" + key + "] must be a decimal value")
	}
}

func (p *Parameters) ValidateDecimalWithInclusiveBounds(key string, value interface{}, minValue float64, maxValue float64) {
	valueAsFloat, typeIsOk := value.(float64)
	if !typeIsOk {
		p.validationErrors.AddMessage("Parameter [" + key + "] must be a decimal value")
		return
	}

	if valueAsFloat < minValue || valueAsFloat > maxValue {
		message := fmt.Sprintf("Parameter [%s] supplied with decimal value [%v], but must be between [%.04f] and [%.04f] inclusive", key, value, minValue, maxValue)
		p.validationErrors.AddMessage(message)
	}
}
