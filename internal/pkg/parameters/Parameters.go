// Copyright (c) 2018 Australian Rivers Institute.

package parameters

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
	"github.com/LindsayBradford/crem/pkg/errors"
)

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

// ContainedLogger is an interface for anything needing Parameters
type Container interface {
	SetParameters(params Map) error
	ParameterErrors() error
}

type Parameters struct {
	paramMap         Map
	specifications   specification.Specifications
	validationErrors errors.CompositeError
}

func (p *Parameters) Enforcing(specs *specification.Specifications) *Parameters {
	p.WithSpecifications(specs).CreatingDefaults()
	return p
}

type Validator func(key string, value interface{}) bool

func (p *Parameters) Initialise(parametersName string) *Parameters {
	p.validationErrors.Initialise(parametersName)
	p.specifications = make(specification.Specifications, 0)
	return p
}

func (p *Parameters) HasEntry(entryKey string) bool {
	_, entryFound := p.paramMap[entryKey]
	return entryFound
}

func (p *Parameters) CreatingDefaults() {
	p.paramMap = make(Map, 0)
	for key, value := range p.specifications {
		if !value.IsOptional {
			p.paramMap[key] = value.DefaultValue
		}
	}
}

func (p *Parameters) WithSpecifications(specifications *specification.Specifications) *Parameters {
	p.specifications = *specifications
	return p
}

func (p *Parameters) AssignAllUserValues(userValues Map) {
	for suppliedKey, suppliedValue := range userValues {
		if p.validateParam(suppliedKey, suppliedValue) {
			p.paramMap[suppliedKey] = suppliedValue
		}
	}
}

func (p *Parameters) AssignOnlyEnforcedUserValues(userValues Map) {
	for _, key := range p.specifications.Keys() {
		if value, userSpecifiedKey := userValues[key]; userSpecifiedKey {
			if p.validateParam(key, value) {
				p.paramMap[key] = value
			}
		}
	}
}

func (p *Parameters) AddValidationErrorMessage(errorMessage string) {
	p.validationErrors.AddMessage(errorMessage)
}

func (p *Parameters) ValidationErrors() error {
	if p.validationErrors.Size() > 0 {
		return &p.validationErrors
	}
	return nil
}

func (p *Parameters) validateParam(key string, value interface{}) bool {
	validationError := p.specifications.Validate(key, value).(specification.ValidationError)
	if !validationError.IsValid() {
		p.validationErrors.Add(validationError)
	}
	return validationError.IsValid()
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

func (p *Parameters) GetBoolean(key string) bool {
	return p.paramMap[key].(bool)
}
