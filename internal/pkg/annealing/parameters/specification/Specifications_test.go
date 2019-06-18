// Copyright (c) 2019 Australian Rivers Institute.

package specification

import (
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
)

const defaultDecimalValue = float64(1)
const defaultIntegerValue = int64(1)
const defaultStringValue = "<undefiled>"

func TestSpecifications_InvalidKey(t *testing.T) {
	g := NewGomegaWithT(t)

	specsUnderTest := New()

	noSpecificationError := specsUnderTest.Validate(notValidKey, defaultIntegerValue).(ValidationError)
	t.Log(noSpecificationError)
	g.Expect(noSpecificationError.IsValid()).To(BeFalse())
}

func TestSpecifications_Decimal(t *testing.T) {
	g := NewGomegaWithT(t)

	specsUnderTest := New()

	specsUnderTest.Add(
		Specification{
			Key:          decimalKey,
			Validator:    IsDecimal,
			DefaultValue: defaultDecimalValue,
		},
	)

	notDecimalError := specsUnderTest.Validate(decimalKey, defaultStringValue).(ValidationError)
	t.Log(notDecimalError)
	g.Expect(notDecimalError.IsValid()).To(BeFalse())

	validError := specsUnderTest.Validate(decimalKey, float64(10)).(ValidationError)
	t.Log(validError)
	g.Expect(validError.IsValid()).To(BeTrue())
}

func TestSpecifications_NonNegativeDecimal(t *testing.T) {
	g := NewGomegaWithT(t)

	specsUnderTest := New()

	specsUnderTest.Add(
		Specification{
			Key:          nonNegativeDecimalKey,
			Validator:    IsNonNegativeDecimal,
			DefaultValue: defaultDecimalValue,
		},
	)

	notDecimalError := specsUnderTest.Validate(nonNegativeDecimalKey, defaultStringValue).(ValidationError)
	t.Log(notDecimalError)
	g.Expect(notDecimalError.IsValid()).To(BeFalse())

	negativeDecimalError := specsUnderTest.Validate(nonNegativeDecimalKey, float64(-10)).(ValidationError)
	t.Log(negativeDecimalError)
	g.Expect(negativeDecimalError.IsValid()).To(BeFalse())

	validError := specsUnderTest.Validate(nonNegativeDecimalKey, float64(10)).(ValidationError)
	t.Log(validError)
	g.Expect(validError.IsValid()).To(BeTrue())
}

func TestSpecifications_BetweenZeroAndOneDecimal(t *testing.T) {
	g := NewGomegaWithT(t)

	specsUnderTest := New()

	specsUnderTest.Add(
		Specification{
			Key:          betweenZeroAndOneDecimalKey,
			Validator:    IsDecimalBetweenZeroAndOne,
			DefaultValue: defaultDecimalValue,
		},
	)

	notDecimalError := specsUnderTest.Validate(betweenZeroAndOneDecimalKey, defaultStringValue).(ValidationError)
	t.Log(notDecimalError)
	g.Expect(notDecimalError.IsValid()).To(BeFalse())

	negativeDecimalError := specsUnderTest.Validate(betweenZeroAndOneDecimalKey, float64(-0.0000000001)).(ValidationError)
	t.Log(negativeDecimalError)
	g.Expect(negativeDecimalError.IsValid()).To(BeFalse())

	validLowerBoundError := specsUnderTest.Validate(betweenZeroAndOneDecimalKey, float64(0)).(ValidationError)
	t.Log(validLowerBoundError)
	g.Expect(validLowerBoundError.IsValid()).To(BeTrue())

	validUpperBoundError := specsUnderTest.Validate(betweenZeroAndOneDecimalKey, float64(1)).(ValidationError)
	t.Log(validUpperBoundError)
	g.Expect(validUpperBoundError.IsValid()).To(BeTrue())

	positiveDecimalError := specsUnderTest.Validate(betweenZeroAndOneDecimalKey, float64(1.0000000001)).(ValidationError)
	t.Log(positiveDecimalError)
	g.Expect(positiveDecimalError.IsValid()).To(BeFalse())
}

func TestSpecifications_Integer(t *testing.T) {
	g := NewGomegaWithT(t)

	specsUnderTest := New()

	specsUnderTest.Add(
		Specification{
			Key:          integerKey,
			Validator:    IsInteger,
			DefaultValue: defaultDecimalValue,
		},
	)

	notIntegerError := specsUnderTest.Validate(integerKey, defaultStringValue).(ValidationError)
	t.Log(notIntegerError)
	g.Expect(notIntegerError.IsValid()).To(BeFalse())

	validError := specsUnderTest.Validate(integerKey, int64(10)).(ValidationError)
	t.Log(validError)
	g.Expect(validError.IsValid()).To(BeTrue())
}

func TestSpecifications_NonNegativeInteger(t *testing.T) {
	g := NewGomegaWithT(t)

	specsUnderTest := New()

	specsUnderTest.Add(
		Specification{
			Key:          nonNegativeIntegerKey,
			Validator:    IsNonNegativeInteger,
			DefaultValue: defaultIntegerValue,
		},
	)

	notIntegerError := specsUnderTest.Validate(nonNegativeIntegerKey, defaultStringValue).(ValidationError)
	t.Log(notIntegerError)
	g.Expect(notIntegerError.IsValid()).To(BeFalse())

	negativeIntegerError := specsUnderTest.Validate(nonNegativeIntegerKey, int64(-10)).(ValidationError)
	t.Log(negativeIntegerError)
	g.Expect(negativeIntegerError.IsValid()).To(BeFalse())

	validError := specsUnderTest.Validate(nonNegativeIntegerKey, int64(10)).(ValidationError)
	t.Log(validError)
	g.Expect(validError.IsValid()).To(BeTrue())
}

func TestSpecifications_IsString(t *testing.T) {
	g := NewGomegaWithT(t)

	specsUnderTest := New()

	specsUnderTest.Add(
		Specification{
			Key:          stringKey,
			Validator:    IsString,
			DefaultValue: defaultStringValue,
		},
	)

	notStringError := specsUnderTest.Validate(stringKey, defaultIntegerValue).(ValidationError)
	t.Log(notStringError)
	g.Expect(notStringError.IsValid()).To(BeFalse())

	validError := specsUnderTest.Validate(stringKey, defaultStringValue).(ValidationError)
	t.Log(validError)
	g.Expect(validError.IsValid()).To(BeTrue())
}

func TestSpecifications_IsReadableFile(t *testing.T) {
	g := NewGomegaWithT(t)

	specsUnderTest := New()

	specsUnderTest.Add(
		Specification{
			Key:          readableFileKey,
			Validator:    IsReadableFile,
			DefaultValue: defaultStringValue,
		},
	)

	notStringError := specsUnderTest.Validate(readableFileKey, defaultIntegerValue).(ValidationError)
	t.Log(notStringError)
	g.Expect(notStringError.IsValid()).To(BeFalse())

	nonExistentFileError := specsUnderTest.Validate(readableFileKey, "noValidFilePresent.txt").(ValidationError)
	t.Log(nonExistentFileError)
	g.Expect(nonExistentFileError.IsValid()).To(BeFalse())

	validError := specsUnderTest.Validate(readableFileKey, "testdata/readableFile.txt").(ValidationError)
	t.Log(validError)
	g.Expect(validError.IsValid()).To(BeTrue())
}

func TestSpecifications_Merge(t *testing.T) {
	g := NewGomegaWithT(t)

	specsUnderTest := New()

	specsUnderTest.Add(
		Specification{
			Key:          decimalKey,
			Validator:    IsDecimal,
			DefaultValue: defaultDecimalValue,
		},
	)

	validError := specsUnderTest.Validate(decimalKey, float64(10)).(ValidationError)
	t.Log(validError)
	g.Expect(validError.IsValid()).To(BeTrue())

	stringSpecMissingError := specsUnderTest.Validate(stringKey, defaultStringValue).(ValidationError)
	t.Log(stringSpecMissingError)
	g.Expect(stringSpecMissingError.IsValid()).To(BeFalse())

	stringSpec := New()

	stringSpec.Add(
		Specification{
			Key:          stringKey,
			Validator:    IsString,
			DefaultValue: defaultStringValue,
		},
	)

	stringSpecPresentError := stringSpec.Validate(stringKey, defaultStringValue).(ValidationError)
	t.Log(stringSpecPresentError)
	g.Expect(stringSpecPresentError.IsValid()).To(BeTrue())

	specsUnderTest.Merge(*stringSpec)

	stringSpecMergedError := specsUnderTest.Validate(stringKey, defaultStringValue).(ValidationError)
	t.Log(stringSpecMergedError)
	g.Expect(stringSpecMergedError.IsValid()).To(BeTrue())
}
