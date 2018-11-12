// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute.

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealers

import (
	"testing"

	"github.com/LindsayBradford/crem/annealing"
	"github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/annealing/parameters"
	"github.com/LindsayBradford/crem/logging"
	"github.com/LindsayBradford/crem/logging/loggers"
	. "github.com/onsi/gomega"
)

type dummyObserver struct{}

func (*dummyObserver) ObserveAnnealingEvent(event annealing.Event) {}

func TestBuild_OverridingDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedLogHandler := new(loggers.BareBonesLogger)
	expectedSolutionExplorer := new(explorer.DumbExplorer)
	expectedObservers := []annealing.Observer{new(dummyObserver)}

	builder := new(Builder)

	expectedParams := parameters.Map {
		StartingTemperature: 1000.0,
		CoolingFactor:       0.5,
		MaximumIterations:   int64(5000),
	}

	annealer, _ := builder.
		SimpleAnnealer().
		WithParameters(expectedParams).
		WithLogHandler(expectedLogHandler).
		WithSolutionExplorer(expectedSolutionExplorer).
		WithObservers(expectedObservers...).
		Build()

	g.Expect(
		annealer.Temperature()).To(BeNumerically("==",expectedParams[StartingTemperature]),
		"Annealer should have built with supplied Temperature")

	g.Expect(
		annealer.CoolingFactor()).To(BeNumerically("==",expectedParams[CoolingFactor]),
		"Annealer should have built with supplied Cooling Factor")

	g.Expect(
		annealer.MaximumIterations()).To(BeNumerically("==",expectedParams[MaximumIterations]),
		"Annealer should have built with supplied Iterations")

	g.Expect(
		annealer.CurrentIteration()).To(BeZero(),
		"Annealer should have built with current iteration of 0")

	g.Expect(
		annealer.LogHandler()).To(BeIdenticalTo(expectedLogHandler),
		"Annealer should have built with supplied Logger")

	g.Expect(
		annealer.SolutionExplorer()).To(BeIdenticalTo(expectedSolutionExplorer),
		"Annealer should have built with supplied Explorer")

	g.Expect(
		annealer.Observers()).To(Equal(expectedObservers),
		"Annealer should have built with supplied Observers")
}

func TestBuild_BadInputs(t *testing.T) {
	g := NewGomegaWithT(t)

	badParams := parameters.Map {
		StartingTemperature: -1,
		CoolingFactor:       1.00000001,
		MaximumIterations:   -1,
	}

	badLogHandler := logging.Logger(nil)
	badExplorer := explorer.Explorer(nil)

	builder := new(Builder)

	annealer, err := builder.
		SimpleAnnealer().
		WithParameters(badParams).
		WithLogHandler(badLogHandler).
		WithSolutionExplorer(badExplorer).
		WithObservers(nil).
		Build()

	g.Expect(
		err.Size()).To(BeNumerically(">", 3),"Annealer should have built with errors")

	t.Log(err)

	g.Expect(
		annealer).To(BeNil(),
		"Annealer should not have been built")
}

func TestAnnealerBuilder_WithDumbSolutionExplorer(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedObjectiveValue := float64(10)

	expectedSolutionExplorer := new(explorer.DumbExplorer)
	expectedSolutionExplorer.SetObjectiveValue(expectedObjectiveValue)
	expectedSolutionExplorer.SetScenarioId("Simple Annealer")

	builder := new(Builder)

	annealer, err := builder.
		SimpleAnnealer().
		WithDumbSolutionExplorer(expectedObjectiveValue).
		Build()

	g.Expect(err).To(BeNil(), "Annealer should have built without errors")

	g.Expect(
		annealer.SolutionExplorer()).To(Equal(expectedSolutionExplorer),
		"Annealer should have built with expected DumbExplorer")

}
