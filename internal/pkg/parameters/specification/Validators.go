// Copyright (c) 2019 Australian Rivers Institute.

package specification

import (
	"fmt"
	"math"
	"os"

	"github.com/pkg/errors"
)

type ValidationError interface {
	error
	IsValid() bool
}

type ValidSpecificationError struct {
	baseError error
	isValid   bool
}

func (m *ValidSpecificationError) Error() string {
	return m.baseError.Error()
}

func (m *ValidSpecificationError) IsValid() bool {
	return m.isValid
}

func NewInvalidSpecificationError(message string) *ValidSpecificationError {
	newError := new(ValidSpecificationError)
	newError.baseError = errors.New(message)
	return newError
}

func NewValidSpecificationError(key string, value interface{}) *ValidSpecificationError {
	newError := new(ValidSpecificationError)
	errorText := fmt.Sprintf("Value [%v] is valid for Parameter [%s]", value, key)
	newError.baseError = errors.New(errorText)
	newError.isValid = true
	return newError
}

func IsDecimal(key string, value interface{}) error {
	_, typeIsOk := value.(float64)
	if !typeIsOk {
		return NewInvalidSpecificationError("Parameter [" + key + "] must be a decimal value")
	}
	return NewValidSpecificationError(key, value)
}

func IsDecimalBetweenZeroAndOne(key string, value interface{}) error {
	return IsDecimalWithInclusiveBounds(key, value, 0, 1)
}

func IsNonNegativeDecimal(key string, value interface{}) error {
	return IsDecimalWithInclusiveBounds(key, value, 0, math.MaxFloat64)
}

func IsDecimalWithInclusiveBounds(key string, value interface{}, minValue float64, maxValue float64) error {
	valueAsFloat, typeIsOk := value.(float64)
	if !typeIsOk {
		return NewInvalidSpecificationError("Parameter [" + key + "] must be a decimal value")
	}

	if valueAsFloat < minValue || valueAsFloat > maxValue {
		message := fmt.Sprintf("Parameter [%s] supplied with decimal value [%g], but must be between [%g] and [%g] inclusive", key, value, minValue, maxValue)
		return NewInvalidSpecificationError(message)
	}
	return NewValidSpecificationError(key, value)
}

func IsInteger(key string, value interface{}) error {
	_, typeIsOk := value.(int64)
	if !typeIsOk {
		return NewInvalidSpecificationError("Parameter [" + key + "] must be an integer value")
	}
	return NewValidSpecificationError(key, value)
}

func IsNonNegativeInteger(key string, value interface{}) error {
	return IsIntegerWithInclusiveBounds(key, value, 0, math.MaxInt64)
}

func IsIntegerWithInclusiveBounds(key string, value interface{}, minValue int64, maxValue int64) error {
	valueAsInteger, typeIsOk := value.(int64)
	if !typeIsOk {
		return NewInvalidSpecificationError("Parameter [" + key + "] must be an integer value")
	}

	if valueAsInteger < minValue || valueAsInteger > maxValue {
		message := fmt.Sprintf("Parameter [%s] supplied with integer value [%v], but must be between [%d] and [%d] inclusive", key, value, minValue, maxValue)
		return NewInvalidSpecificationError(message)
	}
	return NewValidSpecificationError(key, value)
}

func IsString(key string, value interface{}) error {
	_, typeIsOk := value.(string)
	if !typeIsOk {
		return NewInvalidSpecificationError("Parameter [" + key + "] must be a string value")
	}
	return NewValidSpecificationError(key, value)
}

func IsBoolean(key string, value interface{}) error {
	_, typeIsOk := value.(bool)
	if !typeIsOk {
		return NewInvalidSpecificationError("Parameter [" + key + "] must be a boolean value")
	}
	return NewValidSpecificationError(key, value)
}

func IsReadableFile(key string, value interface{}) error {
	valueAsString, typeIsOk := value.(string)
	if !typeIsOk {
		return NewInvalidSpecificationError("Parameter [" + key + "] must be a string")
	}
	if !isReadableFilePath(valueAsString) {
		return NewInvalidSpecificationError("Parameter [" + key + "] must be a valid path to a readable file")
	}
	return NewValidSpecificationError(key, value)
}

func isReadableFilePath(filePath string) bool {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	defer file.Close()

	return err == nil
}
