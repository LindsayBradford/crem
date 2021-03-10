// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/kirkpatrick"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/dumb"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	. "github.com/onsi/gomega"
)

const equalTo = "=="

type dummyObserver struct{}

func (*dummyObserver) ObserveEvent(event observer.Event) {}

func TestBuild_OverridingDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedLogHandler := new(loggers.BareBonesLogger)
	expectedSolutionExplorer := new(null.Explorer)
	expectedObservers := []observer.Observer{new(dummyObserver)}
	expectedId := "someId"
	expectedEventNotifier := new(observer.SynchronousAnnealingEventNotifier)

	builder := new(Builder)

	expectedParams := parameters.Map{
		MaximumIterations: int64(5000),
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

	annealerAttributes := annealer.EventAttributes(observer.FinishedAnnealing)

	actualMaximumIterations := annealerAttributes.Value(MaximumIterations).(uint64)
	g.Expect(actualMaximumIterations).To(BeNumerically(equalTo, expectedParams[MaximumIterations]))

	actualCurrentIteration := annealerAttributes.Value(CurrentIteration).(uint64)
	g.Expect(actualCurrentIteration).To(BeZero())

	g.Expect(annealer.LogHandler()).To(BeIdenticalTo(expectedLogHandler))

	g.Expect(annealer.SolutionExplorer()).To(BeIdenticalTo(expectedSolutionExplorer))

	g.Expect(annealer.Observers()).To(Equal(expectedObservers))

	g.Expect(annealer.EventNotifier()).To(Equal(expectedEventNotifier))

	g.Expect(annealer.Id()).To(Equal(expectedId))
}

func TestBuild_BadInputs(t *testing.T) {
	g := NewGomegaWithT(t)

	badParams := parameters.Map{
		MaximumIterations: -1,
	}

	badLogHandler := logging.Logger(nil)
	badExplorer := explorer.Explorer(nil)
	badEventNotifier := observer.EventNotifier(nil)

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
		err.Size()).To(BeNumerically(">", 3), "Annealer should have built with errors")

	t.Log(err)

	g.Expect(
		annealer).To(BeNil(),
		"Annealer should not have been built")
}

func TestAnnealerBuilder_WithKirkpatrickExplorer(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedSolutionExplorer := kirkpatrick.New()
	expectedSolutionExplorer.WithModel(dumb.NewModel())
	expectedSolutionExplorer.SetId("Simple Annealer")

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
