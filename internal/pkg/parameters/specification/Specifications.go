// Copyright (c) 2019 Australian Rivers Institute.

package specification

import "github.com/pkg/errors"

type SpecValidator func(key string, value interface{}) error

type Specification struct {
	Key          string
	Validator    SpecValidator
	DefaultValue interface{}
	IsOptional   bool
}

func NewSpecifications() *Specifications {
	newSpecs := make(Specifications, 0)
	return &newSpecs
}

type Specifications map[string]Specification

func (s Specifications) Add(spec Specification) Specifications {
	s[spec.Key] = spec
	return s
}

func (s Specifications) HasEntry(key string) bool {
	_, isPresent := s[key]
	return isPresent
}

func (s Specifications) Keys() []string {
	keys := make([]string, 0)
	for key := range s {
		keys = append(keys, key)
	}
	return keys
}

func (s Specifications) Validate(key string, value interface{}) error {
	if s.HasEntry(key) {
		return s[key].Validator(key, value)
	} else {
		return NewSpecificationMissingError(key)
	}
}

type MissingSpecificationError struct {
	baseError error
	isValid   bool
}

func (m *MissingSpecificationError) Error() string {
	return m.baseError.Error()
}

func (m *MissingSpecificationError) IsValid() bool {
	return m.isValid
}

func NewSpecificationMissingError(key string) *MissingSpecificationError {
	newError := new(MissingSpecificationError)
	newError.baseError = errors.New("Parameter [" + key + "] is not supported")
	return newError
}
