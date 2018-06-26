/*
 * Copyright (c) 2018 Australian Rivers Institure. Author: Lindsay Bradford
 */

package annealing

import . "github.com/onsi/gomega"
import "testing"

func TestAnnealerStateFormatWrapper_Defaults(t *testing.T) {
	g := NewGomegaWithT(t)

	const expectedTemperature = "1.000000"
	const expectedCoolingFactor = "1.000000"
	const expectedMaxIterations = "0"
	const expectedCurrentIteration = "0"

	builder := new(AnnealerBuilder)

	annealer, _ := builder.
		SingleObjectiveAnnealer().
		Build()

	wrapperUnderTest := new(AnnealerStateFormatWrapper).Initialise().Wrapping(annealer)

	g.Expect(
		wrapperUnderTest.Temperature()).To(BeIdenticalTo(expectedTemperature),
		"AnnealerStateFormatWrapper should not have defaulted to temperature of \"%s\"", expectedTemperature)

	g.Expect(
		wrapperUnderTest.CoolingFactor()).To(BeIdenticalTo(expectedCoolingFactor),
		"AnnealerStateFormatWrapper should not have defaulted to cooling temperature of \"%s\"", expectedCoolingFactor)

	g.Expect(
		wrapperUnderTest.MaxIterations()).To(BeIdenticalTo(expectedMaxIterations),
		"AnnealerStateFormatWrapper  should not have defaulted to max iterations of \"%s\"", expectedCoolingFactor)

	g.Expect(
		wrapperUnderTest.CurrentIteration()).To(BeIdenticalTo(expectedCurrentIteration),
		"AnnealerStateFormatWrapper  should not have defaulted to current iteration of \"%s\"", expectedCurrentIteration)
}

func TestAnnealerStateFormatWrapper_FormatOverrides(t *testing.T) {
	g := NewGomegaWithT(t)

	const expectedTemperature = "1.00"
	const expectedCoolingFactor = "1.0"
	const expectedMaxIterations = "000"
	const expectedCurrentIteration = "00"

	builder := new(AnnealerBuilder)

	annealer, _ := builder.
		SingleObjectiveAnnealer().
		Build()

	wrapperUnderTest := &AnnealerStateFormatWrapper{}
	wrapperUnderTest.Initialise().Wrapping(annealer)

	wrapperUnderTest.MethodFormats["Temperature"] = "%0.2f"
	g.Expect(
		wrapperUnderTest.Temperature()).To(BeIdenticalTo(expectedTemperature),
		"AnnealerStateFormatWrapper should not have defaulted to temperature of \"%s\"", expectedTemperature)

	wrapperUnderTest.MethodFormats["CoolingFactor"] = "%0.1f"
	g.Expect(
		wrapperUnderTest.CoolingFactor()).To(BeIdenticalTo(expectedCoolingFactor),
		"AnnealerStateFormatWrapper should not have defaulted to cooling temperature of \"%s\"", expectedCoolingFactor)

	wrapperUnderTest.MethodFormats["MaxIterations"] = "%03d"
	g.Expect(
		wrapperUnderTest.MaxIterations()).To(BeIdenticalTo(expectedMaxIterations),
		"AnnealerStateFormatWrapper  should not have defaulted to max iterations of \"%s\"", expectedCoolingFactor)

	wrapperUnderTest.MethodFormats["CurrentIteration"] = "%02d"
	g.Expect(
		wrapperUnderTest.CurrentIteration()).To(BeIdenticalTo(expectedCurrentIteration),
		"AnnealerStateFormatWrapper  should not have defaulted to current iteration of \"%s\"", expectedCurrentIteration)
}