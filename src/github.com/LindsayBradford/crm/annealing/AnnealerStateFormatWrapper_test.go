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

	annealer := builder.
		SingleObjectiveAnnealer().
		Build()

	wrapperUnderTest := NewAnnealerStateFormatWrapper(annealer)

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
