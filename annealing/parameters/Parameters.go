// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute.

package parameters

import (
	"fmt"

	"github.com/LindsayBradford/crem/errors"
)

type Map map[string]interface{}
type MetaDataMap map[string]MetaData

type Parameters struct {
	paramMap         Map
	metaDataMap      MetaDataMap
	validationErrors *errors.CompositeError
}

type Validator func(key string, value interface{}) bool

type MetaData struct {
	Key          string
	Validator    Validator
	DefaultValue interface{}
}

func (p *Parameters) Initialise() *Parameters {
	p.validationErrors = errors.NewComposite("SolutionExplorer Parameters")
	p.metaDataMap = make(MetaDataMap, 0)
	return p
}

func (p *Parameters) CreateDefaults() {
	p.paramMap = make(Map, 0)
	for key, value := range p.metaDataMap {
		p.paramMap[key] = value.DefaultValue
	}
}

func (p *Parameters) AddMetaData(metaData MetaData) {
	p.metaDataMap[metaData.Key] = metaData
}

func (p *Parameters) Merge(params Map) {
	p.validationErrors = errors.NewComposite("SolutionExplorer Parameters")
	for suppliedKey, suppliedValue := range params {
		if p.validateParam(suppliedKey, suppliedValue) {
			p.paramMap[suppliedKey] = suppliedValue
		}
	}
}

func (p *Parameters) ValidationErrors() error {
	if p.validationErrors.Size() > 0 {
		return p.validationErrors
	}
	return nil
}

func (p *Parameters) validateParam(key string, value interface{}) bool {
	if _, isPresent := p.metaDataMap[key]; isPresent {
		return p.metaDataMap[key].Validator(key, value)
	} else {
		p.keyMissingValidator(key)
	}
	return true
}

func (p *Parameters) keyMissingValidator(key string) {
	p.validationErrors.AddMessage("Parameter [" + string(key) + "] is not a parameter for this explorer")
}

func (p *Parameters) ValidateIsDecimal(key string, value interface{}) bool {
	_, typeIsOk := value.(float64)
	if !typeIsOk {
		p.validationErrors.AddMessage("Parameter [" + key + "] must be a decimal value")
		return false
	}
	return true
}

func (p *Parameters) ValidateDecimalWithInclusiveBounds(key string, value interface{}, minValue float64, maxValue float64) bool {
	valueAsFloat, typeIsOk := value.(float64)
	if !typeIsOk {
		p.validationErrors.AddMessage("Parameter [" + key + "] must be a decimal value")
		return false
	}

	if valueAsFloat < minValue || valueAsFloat > maxValue {
		message := fmt.Sprintf("Parameter [%s] supplied with decimal value [%v], but must be between [%.04f] and [%.04f] inclusive", key, value, minValue, maxValue)
		p.validationErrors.AddMessage(message)
		return false
	}
	return true
}

func (p *Parameters) Get(key string) interface{} {
	return p.paramMap[key]
}
