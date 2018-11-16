// Copyright (c) 2018 Australian Rivers Institute.

package parameters

import (
	"fmt"
	"math"
	"os"

	"github.com/LindsayBradford/crem/pkg/errors"
)

type Parameters struct {
	paramMap         Map
	metaDataMap      MetaDataMap
	validationErrors *errors.CompositeError
}

type Map map[string]interface{}

func (m Map) SetInt64(key string, value int64) {
	m[key] = value
}

func (m Map) SetFloat64(key string, value float64) {
	m[key] = value
}

func (m Map) SetString(key string, value string) {
	m[key] = value
}

type MetaDataMap map[string]MetaData

type MetaData struct {
	Key          string
	Validator    Validator
	DefaultValue interface{}
}

type Validator func(key string, value interface{}) bool

func (p *Parameters) Initialise() *Parameters {
	p.validationErrors = errors.New("SolutionExplorer Parameters")
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
	p.validationErrors = errors.New("SolutionExplorer Parameters")
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

func (p *Parameters) IsDecimal(key string, value interface{}) bool {
	_, typeIsOk := value.(float64)
	if !typeIsOk {
		p.validationErrors.AddMessage("Parameter [" + key + "] must be a decimal value")
		return false
	}
	return true
}

func (p *Parameters) IsDecimalBetweenZeroAndOne(key string, value interface{}) bool {
	return p.IsDecimalWithInclusiveBounds(key, value, 0, 1)
}

func (p *Parameters) IsNonNegativeDecimal(key string, value interface{}) bool {
	return p.IsDecimalWithInclusiveBounds(key, value, 0, math.MaxFloat64)
}

func (p *Parameters) IsDecimalWithInclusiveBounds(key string, value interface{}, minValue float64, maxValue float64) bool {
	valueAsFloat, typeIsOk := value.(float64)
	if !typeIsOk {
		p.validationErrors.AddMessage("Parameter [" + key + "] must be a decimal value")
		return false
	}

	if valueAsFloat < minValue || valueAsFloat > maxValue {
		message := fmt.Sprintf("Parameter [%s] supplied with decimal value [%g], but must be between [%g] and [%g] inclusive", key, value, minValue, maxValue)
		p.validationErrors.AddMessage(message)
		return false
	}
	return true
}

func (p *Parameters) IsInteger(key string, value interface{}) bool {
	_, typeIsOk := value.(int64)
	if !typeIsOk {
		p.validationErrors.AddMessage("Parameter [" + key + "] must be an integer value")
		return false
	}
	return true
}

func (p *Parameters) IsNonNegativeInteger(key string, value interface{}) bool {
	return p.IsIntegerWithInclusiveBounds(key, value, 0, math.MaxInt64)
}

func (p *Parameters) IsIntegerWithInclusiveBounds(key string, value interface{}, minValue int64, maxValue int64) bool {
	valueAsInteger, typeIsOk := value.(int64)
	if !typeIsOk {
		p.validationErrors.AddMessage("Parameter [" + key + "] must be an integer value")
		return false
	}

	if valueAsInteger < minValue || valueAsInteger > maxValue {
		message := fmt.Sprintf("Parameter [%s] supplied with integer value [%v], but must be between [%d] and [%d] inclusive", key, value, minValue, maxValue)
		p.validationErrors.AddMessage(message)
		return false
	}
	return true
}

func (p *Parameters) IsReadableFile(key string, value interface{}) bool {
	valueAsString, typeIsOk := value.(string)
	if !typeIsOk {
		p.validationErrors.AddMessage("Parameter [" + key + "] must be a string")
		return false
	}
	if !isReadableFilePath(valueAsString) {
		p.validationErrors.AddMessage("Parameter [" + key + "] must be a valid path to a readable file")
		return false
	}
	return true
}

func isReadableFilePath(filePath string) bool {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	defer file.Close()

	return err == nil
}

func (p *Parameters) GetInt64(key string) int64 {
	return p.paramMap[key].(int64)
}

func (p *Parameters) GetFloat64(key string) float64 {
	return p.paramMap[key].(float64)
}

func (p *Parameters) GetString(key string) string {
	return p.paramMap[key].(string)
}
