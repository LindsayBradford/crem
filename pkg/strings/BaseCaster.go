// Copyright (c) 2019 Australian Rivers Institute.

package strings

import (
	"strconv"
)

// BaseCaster allows strings to be cast into other base types where possible.  With a fluent interface, it allows
// the caller fine-tuned control over how strings are cast to other base-types.
type BaseCaster struct {
	numbersAsFloats bool
}

// WithNumbersAsFloats attempts to convert all numbers to floating point number base type.
func (bs *BaseCaster) WithNumbersAsFloats() *BaseCaster {
	bs.numbersAsFloats = true
	return bs
}

// Cast attempts to convert the value supplied as a string to another Golang base type, depending on how the BaseCaster
// has been configured.  If it cannot, the value is returned with its original string type.
func (bs *BaseCaster) Cast(valueAsString string) interface{} {
	valueAsNumber := bs.attemptCastToNumber(valueAsString)
	if castWorked(valueAsNumber) {
		return valueAsNumber
	}
	valueAsBoolean := attemptCastToBoolean(valueAsString)
	if castWorked(valueAsBoolean) {
		return valueAsBoolean
	}
	return valueAsString
}

func (bs *BaseCaster) attemptCastToNumber(value string) interface{} {
	if !bs.numbersAsFloats {
		valueAsUInt, uintError := strconv.ParseUint(value, 10, 64)
		if uintError == nil {
			return valueAsUInt
		}
		valueAsInt, intError := strconv.ParseInt(value, 10, 64)
		if intError == nil {
			return valueAsInt
		}
	}
	valueAsFloat, floatError := strconv.ParseFloat(value, 64)
	if floatError == nil {
		return valueAsFloat
	}

	return value
}

func attemptCastToBoolean(value string) interface{} {
	valueAsBool, boolError := strconv.ParseBool(value)
	if boolError == nil {
		return valueAsBool
	}
	return value
}

func castWorked(value interface{}) bool {
	_, isString := value.(string)
	return !isString
}
