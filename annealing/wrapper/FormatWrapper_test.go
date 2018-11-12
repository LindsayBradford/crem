// Copyright (c) 2018 Australian Rivers Institute.

package wrapper

import (
	"testing"

	"github.com/LindsayBradford/crem/annealing/annealers"
	. "github.com/onsi/gomega"
)

func TestFormatWrapper_Defaults(t *testing.T) {
	g := NewGomegaWithT(t)

	const expectedTemperature = "1.000000"
	const expectedCoolingFactor = "1.000000"
	const expectedMaximumIterations = "0"
	const expectedCurrentIteration = "0"

	annealer := new(annealers.SimpleAnnealer)
	annealer.Initialise()

	wrapperUnderTest := new(FormatWrapper).Initialise().Wrapping(annealer)

	g.Expect(
		wrapperUnderTest.Temperature()).To(BeIdenticalTo(expectedTemperature),
		"FormatWrapper should not have defaulted to temperature of \"%s\"", expectedTemperature)

	g.Expect(
		wrapperUnderTest.CoolingFactor()).To(BeIdenticalTo(expectedCoolingFactor),
		"FormatWrapper should not have defaulted to cooling temperature of \"%s\"", expectedCoolingFactor)

	g.Expect(
		wrapperUnderTest.MaximumIterations()).To(BeIdenticalTo(expectedMaximumIterations),
		"FormatWrapper  should not have defaulted to max iterations of \"%s\"", expectedCoolingFactor)

	g.Expect(
		wrapperUnderTest.CurrentIteration()).To(BeIdenticalTo(expectedCurrentIteration),
		"FormatWrapper  should not have defaulted to current iteration of \"%s\"", expectedCurrentIteration)
}

func TestFormatWrapper_FormatOverrides(t *testing.T) {
	g := NewGomegaWithT(t)

	const expectedTemperature = "1.00"
	const expectedCoolingFactor = "1.0"
	const expectedMaximumIterations = "000"
	const expectedCurrentIteration = "00"

	annealer := new(annealers.SimpleAnnealer)
	annealer.Initialise()

	wrapperUnderTest := &FormatWrapper{}
	wrapperUnderTest.Initialise().Wrapping(annealer)

	wrapperUnderTest.MethodFormats["Temperature"] = "%0.2f"
	g.Expect(
		wrapperUnderTest.Temperature()).To(BeIdenticalTo(expectedTemperature),
		"FormatWrapper should not have defaulted to temperature of \"%s\"", expectedTemperature)

	wrapperUnderTest.MethodFormats["CoolingFactor"] = "%0.1f"
	g.Expect(
		wrapperUnderTest.CoolingFactor()).To(BeIdenticalTo(expectedCoolingFactor),
		"FormatWrapper should not have defaulted to cooling temperature of \"%s\"", expectedCoolingFactor)

	wrapperUnderTest.MethodFormats["MaximumIterations"] = "%03d"
	g.Expect(
		wrapperUnderTest.MaximumIterations()).To(BeIdenticalTo(expectedMaximumIterations),
		"FormatWrapper  should not have defaulted to max iterations of \"%s\"", expectedCoolingFactor)

	wrapperUnderTest.MethodFormats["CurrentIteration"] = "%02d"
	g.Expect(
		wrapperUnderTest.CurrentIteration()).To(BeIdenticalTo(expectedCurrentIteration),
		"FormatWrapper  should not have defaulted to current iteration of \"%s\"", expectedCurrentIteration)
}
