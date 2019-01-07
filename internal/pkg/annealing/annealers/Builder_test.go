// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute.

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealers

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/kirkpatrick"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/dumb"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	. "github.com/onsi/gomega"
)

type dummyObserver struct{}

func (*dummyObserver) ObserveAnnealingEvent(event annealing.Event) {}

func TestBuild_OverridingDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedLogHandler := new(loggers.BareBonesLogger)
	expectedSolutionExplorer := new(null.Explorer)
	expectedObservers := []annealing.Observer{new(dummyObserver)}
	expectedId := "someId"
	expectedEventNotifier := new(annealing.ConcurrentAnnealingEventNotifier)

	builder := new(Builder)

	expectedParams := parameters.Map{
		StartingTemperature: 1000.0,
		CoolingFactor:       0.5,
		MaximumIterations:   int64(5000),
	}

	annealer, _ := builder.
		SimpleAnnealer().
		WithId(expectedId).
		WithParameters(expectedParams).
		WithLogHandler(expectedLogHandler).
		WithSolutionExplorer(expectedSolutionExplorer).
		WithEventNotifier(expectedEventNotifier).
		WithObservers(expectedObservers...).
		Build()

	g.Expect(
		annealer.Temperature()).To(BeNumerically("==", expectedParams[StartingTemperature]),
		"Annealer should have built with supplied Temperature")

	g.Expect(
		annealer.CoolingFactor()).To(BeNumerically("==", expectedParams[CoolingFactor]),
		"Annealer should have built with supplied Cooling Factor")

	g.Expect(
		annealer.MaximumIterations()).To(BeNumerically("==", expectedParams[MaximumIterations]),
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

	g.Expect(
		annealer.EventNotifier()).To(Equal(expectedEventNotifier),
		"Annealer should have built with supplied EventNotifier")

	g.Expect(
		annealer.Id()).To(Equal(expectedId),
		"Annealer should have built with supplied Id")
}

func TestBuild_BadInputs(t *testing.T) {
	g := NewGomegaWithT(t)

	badParams := parameters.Map{
		StartingTemperature: -1,
		CoolingFactor:       1.00000001,
		MaximumIterations:   -1,
	}

	badLogHandler := logging.Logger(nil)
	badExplorer := explorer.Explorer(nil)
	badEventNotifier := annealing.EventNotifier(nil)

	builder := new(Builder)

	annealer, err := builder.
		SimpleAnnealer().
		WithParameters(badParams).
		WithLogHandler(badLogHandler).
		WithSolutionExplorer(badExplorer).
		WithEventNotifier(badEventNotifier).
		WithObservers(nil).
		Build()

	g.Expect(
		err.Size()).To(BeNumerically(">", 4), "Annealer should have built with errors")

	t.Log(err)

	g.Expect(
		annealer).To(BeNil(),
		"Annealer should not have been built")
}

func TestAnnealerBuilder_WithKirkpatrickExplorer(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedSolutionExplorer := kirkpatrick.New()
	expectedSolutionExplorer.WithModel(dumb.New())
	expectedSolutionExplorer.SetScenarioId("Simple Annealer")

	builder := new(Builder)

	annealer, err := builder.
		SimpleAnnealer().
		WithDumbSolutionExplorer().
		Build()

	g.Expect(err).To(BeNil(), "Annealer should have built without errors")

	g.Expect(
		annealer.SolutionExplorer()).To(BeAssignableToTypeOf(expectedSolutionExplorer),
		"Annealer should have built with expected DumbExplorer")

}
