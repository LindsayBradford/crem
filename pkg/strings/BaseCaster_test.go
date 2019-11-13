// Copyright (c) 2019 Australian Rivers Institute.

package strings

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBaseCaster_NewBaseCaster(t *testing.T) {
	g := NewGomegaWithT(t)

	casterUnderTest := new(BaseCaster)

	g.Expect(*casterUnderTest).To(BeAssignableToTypeOf(BaseCaster{}))
}

func TestBaseCaster_Cast(t *testing.T) {
	g := NewGomegaWithT(t)

	casterUnderTest := new(BaseCaster)

	verifyCastTypeIsDesired := func(value string, valueWithDesiredType interface{}) {
		actualValue := casterUnderTest.Cast(value)
		g.Expect(actualValue).To(BeAssignableToTypeOf(valueWithDesiredType))
	}

	verifyCastTypeIsDesired("This can be no other base type than a string", "isString")
	verifyCastTypeIsDesired("42", uint64(42))
	verifyCastTypeIsDesired("-42", int64(-42))
	verifyCastTypeIsDesired("-42.42", float64(-42.42))
	verifyCastTypeIsDesired("true", false)
}

func TestBaseCaster_Cast_WithNumbersAsFloats(t *testing.T) {
	g := NewGomegaWithT(t)

	casterUnderTest := new(BaseCaster).WithNumbersAsFloats()

	verifyCastTypeIsFloat64 := func(value string) {
		const valueWithExpectedType = float64(42)
		actualValue := casterUnderTest.Cast(value)
		g.Expect(actualValue).To(BeAssignableToTypeOf(valueWithExpectedType))
	}

	verifyCastTypeIsFloat64("42")
	verifyCastTypeIsFloat64("-42")
}

func ExampleBaseCaster_Cast() {
	caster := new(BaseCaster).WithNumbersAsFloats()

	fortyTwoAsString := "42"
	sevenAsString := "7"

	fortyTwo, fortyTwoIsFloat := caster.Cast(fortyTwoAsString).(float64)
	seven, sevenIsFloat := caster.Cast(sevenAsString).(float64)

	var actualResult float64
	if sevenIsFloat && fortyTwoIsFloat {
		actualResult = fortyTwo / seven
	}

	fmt.Printf("%v", actualResult)

	// Output: 6
}
