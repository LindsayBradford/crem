// Copyright (c) 2018 Australian Rivers Institute.

package parameters

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"
)

func TestEmptyParameters_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()

	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "No errors on initialising empty parameters")
}

func TestParametersAddMetaData_CreateDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	g.Expect(parametersUnderTest.Get("testKey")).To(BeNumerically("==", 0.69), "metadata should have set correct default")
}

func TestParameters_MergeValidParameter_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)
	paramsToMerge["testKey"] = 0.5

	parametersUnderTest.Merge(paramsToMerge)

	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "No errors on initialising empty parameters")

	g.Expect(parametersUnderTest.Get("testKey")).To(BeNumerically("==", 0.5), "metadata should have set correct default")
}

func TestParameters_MergeUnknownParameter_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)
	paramsToMerge["notAKnownKey"] = 0.5

	parametersUnderTest.Merge(paramsToMerge)

	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging unknown parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	g.Expect(parametersUnderTest.Get("testKey")).To(BeNumerically("==", 0.69), "metadata should have set correct default")
}

func TestParameters_MergeInvalidParameter_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)
	paramsToMerge["testKey"] = "DefinitelyNotADecimal"

	parametersUnderTest.Merge(paramsToMerge)

	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging unknown parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	g.Expect(parametersUnderTest.Get("testKey")).To(BeNumerically("==", 0.69), "metadata should have set correct default")
}

func TestParameters_MergeInvalidBoundParameter_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	parametersUnderTest := new(Parameters).Initialise()
	addMetaDataUnderTest(parametersUnderTest)

	parametersUnderTest.CreateDefaults()

	paramsToMerge := make(Map, 0)

	paramsToMerge["testBoundKey"] = 0.0004
	parametersUnderTest.Merge(paramsToMerge)

	g.Expect(parametersUnderTest.ValidationErrors()).To(BeNil(), "no error on merging valid parameter")

	paramsToMerge["testBoundKey"] = 0.00051
	parametersUnderTest.Merge(paramsToMerge)

	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging invalid parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	paramsToMerge["testBoundKey"] = 0.999
	parametersUnderTest.Merge(paramsToMerge)

	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging invalid parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	paramsToMerge["testBoundKey"] = "NotEvenADecimal"
	parametersUnderTest.Merge(paramsToMerge)

	g.Expect(parametersUnderTest.ValidationErrors()).To(Not(BeNil()), "error on merging invalid parameter")
	t.Log(parametersUnderTest.ValidationErrors())

	g.Expect(parametersUnderTest.Get("testBoundKey")).To(BeNumerically("==", 0.0004), "metadata should have set correct default")
}

func addMetaDataUnderTest(params *Parameters) {
	params.AddMetaData(
		MetaData{
			Key:          "testKey",
			Validator:    params.ValidateIsDecimal,
			DefaultValue: 0.69,
		},
	)

	validateIsBankErosionFudgeFactor := func(key string, value interface{}) bool {
		minValue := 1 * math.Pow(10, -4)
		maxValue := 5 * math.Pow(10, -4)
		return params.ValidateDecimalWithInclusiveBounds(key, value, minValue, maxValue)
	}

	params.AddMetaData(
		MetaData{
			Key:          "testBoundKey",
			Validator:    validateIsBankErosionFudgeFactor,
			DefaultValue: 0.0005,
		},
	)

}
