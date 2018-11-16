// Copyright (c) 2018 Australian Rivers Institute.

package parameters

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

	readableFileKey = "readableFileKey"
)

const defaultDecimalValue = float64(1)
const defaultIntegerValue = int64(1)
const defaultStringValue = "<undefiled>"

func TestEmptyParameters_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()

	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "No errors on initialising empty parameters")
}

func TestAddMetaData_CreateDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	g.Expect(parametersUnderTest.GetFloat64(decimalKey)).To(BeNumerically("==", defaultDecimalValue), "metadata should have set correct default")
	g.Expect(parametersUnderTest.GetInt64(integerKey)).To(BeNumerically("==", defaultIntegerValue), "metadata should have set correct default")
	g.Expect(parametersUnderTest.GetString(readableFileKey)).To(Equal(defaultStringValue), "metadata should have set correct default")
}

func TestMergeValidDecimal_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge.SetFloat64(decimalKey, -0.5)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "No errors on initialising empty parameters")
	g.Expect(parametersUnderTest.GetFloat64(decimalKey)).To(BeNumerically("==", -0.5), "metadata should have set correct default")

	paramsToMerge.SetFloat64(decimalKey, 0.5)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "No errors on initialising empty parameters")
	g.Expect(parametersUnderTest.GetFloat64(decimalKey)).To(BeNumerically("==", 0.5), "metadata should have set correct default")
}

func TestMergeValidInteger_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge.SetInt64(integerKey, -5)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "No errors on initialising empty parameters")
	g.Expect(parametersUnderTest.GetInt64(integerKey)).To(BeNumerically("==", -5), "metadata should have set correct default")

	paramsToMerge.SetInt64(integerKey, 5)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "No errors on initialising empty parameters")
	g.Expect(parametersUnderTest.GetInt64(integerKey)).To(BeNumerically("==", 5), "metadata should have set correct default")
}

func TestMergeValidFilePath_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge.SetString(readableFileKey, "testdata/readableFile.txt")
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "No errors on initialising empty parameters")
	g.Expect(parametersUnderTest.GetString(readableFileKey)).To(Equal("testdata/readableFile.txt"), "metadata should have set correct default")
}

func TestMergeUnknownParameter_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge.SetFloat64(notValidKey, 0.5)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging unknown parameter")
	t.Log(parametersUnderTest.ValidationErrors())
}

func TestMergeInvalidDecimal_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge[decimalKey] = "DefinitelyNotADecimal"
	parametersUnderTest.Merge(paramsToMerge)

	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging unknown parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	g.Expect(parametersUnderTest.GetFloat64(decimalKey)).To(BeNumerically("==", defaultDecimalValue), "metadata should have set correct default")
}

func TestMergeInvalidInteger_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge[integerKey] = "DefinitelyNotAnInteger"
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging unknown parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	g.Expect(parametersUnderTest.GetInt64(integerKey)).To(BeNumerically("==", defaultIntegerValue), "metadata should have set correct default")
}

func TestMergeInvalidNonNegativeDecimals_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge.SetFloat64(nonNegativeDecimalKey, -0.00001)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging invalid parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	paramsToMerge[nonNegativeDecimalKey] = "NotEvenADecimal"
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging invalid parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	g.Expect(parametersUnderTest.GetFloat64(nonNegativeDecimalKey)).To(BeNumerically("==", defaultDecimalValue), "metadata should have set correct default")
}

func TestMergeBetweenZeroAndOneDecimal(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	validLowerBoundDecimalValue := float64(0)
	validUpperBoundDecimalValue := float64(1)

	paramsToMerge.SetFloat64(betweenZeroAndOneDecimalKey, validLowerBoundDecimalValue)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "no error on merging valid parameter")

	paramsToMerge.SetFloat64(betweenZeroAndOneDecimalKey, validUpperBoundDecimalValue)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "no error on merging valid parameter")

	paramsToMerge.SetFloat64(betweenZeroAndOneDecimalKey, -0.00001)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging invalid parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	paramsToMerge.SetFloat64(betweenZeroAndOneDecimalKey, 1.00001)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging invalid parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	paramsToMerge[betweenZeroAndOneDecimalKey] = "NotEvenADecimal"
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging invalid parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	g.Expect(parametersUnderTest.GetFloat64(betweenZeroAndOneDecimalKey)).To(BeNumerically("==", validUpperBoundDecimalValue), "metadata should have set correct default")
}

func TestMergeValidNonNegativeInteger_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge.SetInt64(nonNegativeIntegerKey, 5)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "No errors on initialising empty parameters")
	g.Expect(parametersUnderTest.GetInt64(nonNegativeIntegerKey)).To(BeNumerically("==", 5), "metadata should have set correct default")
}

func TestMergeInvalidNonNegativeIntegers_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge.SetInt64(nonNegativeIntegerKey, -1)
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging invalid parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	paramsToMerge[nonNegativeIntegerKey] = "NotEvenAnInteger"
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging invalid parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	g.Expect(parametersUnderTest.GetInt64(nonNegativeIntegerKey)).To(BeNumerically("==", defaultIntegerValue), "metadata should have set correct default")
}

func TestMergeInvalidFilePath_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge[readableFileKey] = 0.4
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging unknown parameter")
	t.Log(parametersUnderTest.ValidationErrors())


	paramsToMerge.SetString(readableFileKey, "a non-existent file path")
	parametersUnderTest.Merge(paramsToMerge)
	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging unknown parameter")
	t.Log(parametersUnderTest.ValidationErrors())
}


func addMetaDataUnderTest(params *Parameters) {
	params.AddMetaData(
		MetaData{
			Key:          decimalKey,
			Validator:    params.IsDecimal,
			DefaultValue: defaultDecimalValue,
		},
	)

	params.AddMetaData(
		MetaData{
			Key:          nonNegativeDecimalKey,
			Validator:    params.IsNonNegativeDecimal,
			DefaultValue: defaultDecimalValue,
		},
	)

	params.AddMetaData(
		MetaData{
			Key:          betweenZeroAndOneDecimalKey,
			Validator:    params.IsDecimalBetweenZeroAndOne,
			DefaultValue: defaultDecimalValue,
		},
	)

	params.AddMetaData(
		MetaData{
			Key:          integerKey,
			Validator:    params.IsInteger,
			DefaultValue: defaultIntegerValue,
		},
	)

	params.AddMetaData(
		MetaData{
			Key:          nonNegativeIntegerKey,
			Validator:    params.IsNonNegativeInteger,
			DefaultValue: defaultIntegerValue,
		},
	)

	params.AddMetaData(
		MetaData{
			Key:          readableFileKey,
			Validator:    params.IsReadableFile,
			DefaultValue: defaultStringValue,
		},
	)
}
