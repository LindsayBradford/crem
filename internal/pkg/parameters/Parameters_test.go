// Copyright (c) 2018 Australian Rivers Institute.

package parameters

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	_ = iota

	notValidKey = "notValidKey"

	decimalKey                  = "decimalKey"
	nonNegativeDecimalKey       = "nonNegativeDecimalKey"
	betweenZeroAndOneDecimalKey = "betweenZeroAndOneDecimalKey"

	integerKey            = "integerKey"
	nonNegativeIntegerKey = "nonNegativeIntegerKey"

	stringKey       = "stringKey"
	readableFileKey = "readableFileKey"

	optionalKey = "optionalKey"
)

const defaultDecimalValue = float64(1)
const defaultIntegerValue = int64(1)
const defaultStringValue = "<undefiled>"

func TestEmptyParameters_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise("test")

	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "No errors on initialising empty parameters")
}

func TestAddMetaData_CreateDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters)
	parametersUnderTest.Initialise("test").Enforcing(testSpecifications())

	g.Expect(parametersUnderTest.GetFloat64(decimalKey)).To(BeNumerically("==", defaultDecimalValue), "metadata should have set correct default")
	g.Expect(parametersUnderTest.GetInt64(integerKey)).To(BeNumerically("==", defaultIntegerValue), "metadata should have set correct default")
	g.Expect(parametersUnderTest.GetString(stringKey)).To(Equal(defaultStringValue), "metadata should have set correct default")
	g.Expect(parametersUnderTest.GetString(readableFileKey)).To(Equal(defaultStringValue), "metadata should have set correct default")
	g.Expect(parametersUnderTest.HasEntry(optionalKey)).To(BeFalse(), "default for optional meta data entry should not have been created")
}

func TestMergeOptionalParameter(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters)
	parametersUnderTest.Initialise("test").Enforcing(testSpecifications())

	g.Expect(parametersUnderTest.HasEntry(optionalKey)).To(BeFalse(), "default for optional meta data entry should not have been created")

	paramsToMerge := make(Map, 0)

	expectedValue := 0.4

	paramsToMerge[optionalKey] = expectedValue
	parametersUnderTest.AssignAllUserValues(paramsToMerge)

	g.Expect(parametersUnderTest.GetFloat64(optionalKey)).To(BeNumerically("==", expectedValue), "metadata should have set correct default")
}

func TestAddValidationErrorMessage(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise("test")
	parametersUnderTest.AddValidationErrorMessage("here is a user-defined validation error, useful for embedding semantics tests to one or more parameters")

	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()))
	t.Log(parametersUnderTest.ValidationErrors())
}

func testSpecifications() *specification.Specifications {
	specs := specification.NewSpecifications()
	specs.Add(
		specification.Specification{
			Key:          decimalKey,
			Validator:    specification.IsDecimal,
			DefaultValue: defaultDecimalValue,
		},
	).Add(
		specification.Specification{
			Key:          nonNegativeDecimalKey,
			Validator:    specification.IsNonNegativeDecimal,
			DefaultValue: defaultDecimalValue,
		},
	).Add(
		specification.Specification{
			Key:          betweenZeroAndOneDecimalKey,
			Validator:    specification.IsDecimalBetweenZeroAndOne,
			DefaultValue: defaultDecimalValue,
		},
	).Add(
		specification.Specification{
			Key:          integerKey,
			Validator:    specification.IsInteger,
			DefaultValue: defaultIntegerValue,
		},
	).Add(
		specification.Specification{
			Key:          nonNegativeIntegerKey,
			Validator:    specification.IsNonNegativeInteger,
			DefaultValue: defaultIntegerValue,
		},
	).Add(
		specification.Specification{
			Key:          stringKey,
			Validator:    specification.IsString,
			DefaultValue: defaultStringValue,
		},
	).Add(
		specification.Specification{
			Key:          readableFileKey,
			Validator:    specification.IsReadableFile,
			DefaultValue: defaultStringValue,
		},
	).Add(
		specification.Specification{
			Key:        optionalKey,
			Validator:  specification.IsDecimal,
			IsOptional: true,
		},
	)
	return specs
}
