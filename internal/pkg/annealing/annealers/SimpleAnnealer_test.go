// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"errors"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	. "github.com/onsi/gomega"
)

func TestSimpleAnnealer_Initialise(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	annealerAttributes := annealer.EventAttributes(observer.FinishedAnnealing)

	actualMaximumIterations := annealerAttributes.Value(MaximumIterations).(uint64)
	g.Expect(actualMaximumIterations).To(BeZero())

	actualCurrentIteration := annealerAttributes.Value(CurrentIteration).(uint64)
	g.Expect(actualCurrentIteration).To(BeZero())

	g.Expect(annealer.LogHandler()).To(Equal(loggers.NewNullLogger()))

	g.Expect(annealer.SolutionExplorer()).To(Equal(null.NullExplorer))

	g.Expect(annealer.Observers()).To(BeNil())
}

func TestSimpleAnnealer_DeepClone(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	actualClone := annealer.DeepClone()

	g.Expect(actualClone).To(Equal(annealer))
}

func TestSimpleAnnealer_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	annealer.SetLogHandler(nil)

	g.Expect(annealer.LogHandler()).To(Equal(loggers.NewNullLogger()),
		"Annealer should have forced nil Logger set attempt to DefaultNullLogger")

	explorerErr := annealer.SetSolutionExplorer(nil)

	g.Expect(explorerErr).To(Not(BeNil()))
	g.Expect(annealer.SolutionExplorer()).To(Equal(null.NullExplorer),
		"Annealer should have ignored crap Solution Explorer set attempt")

	observersErr := annealer.AddObserver(nil)

	g.Expect(observersErr).To(Not(BeNil()))
	g.Expect(annealer.Observers()).To(BeNil(),
		"Annealer should have ignored crap AnnealerObserver set attempt")
}

func TestSimpleAnnealer_Anneal(t *testing.T) {
	g := NewGomegaWithT(t)

	const iterations uint64 = 3

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	expectedParams := parameters.Map{
		MaximumIterations: int64(iterations),
	}

	annealer.SetParameters(expectedParams)

	beforeAttributes := annealer.EventAttributes(observer.FinishedAnnealing)
	actualBeforeCurrentIteration := beforeAttributes.Value(CurrentIteration).(uint64)
	g.Expect(actualBeforeCurrentIteration).To(BeZero())

	annealer.Anneal()

	afterAttributes := annealer.EventAttributes(observer.FinishedAnnealing)
	actualAfterCurrentIteration := afterAttributes.Value(CurrentIteration).(uint64)

	g.Expect(actualAfterCurrentIteration).To(BeNumerically("==", iterations))
}

func TestSimpleAnnealer_AddObserver(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	const expectedIterations = uint64(3)

	expectedParams := parameters.Map{
		MaximumIterations: int64(expectedIterations),
	}

	annealer.SetParameters(expectedParams)

	g.Expect(annealer.Observers()).To(BeNil())

	observerError := annealer.AddObserver(nil)

	g.Expect(observerError).To(Not(BeNil()), "Annealer should have raised an error on adding nil AnnealerObserver")

	counter := new(CountingObserver)
	counter.eventCounts = make(map[observer.EventType]uint64)

	observerError = annealer.AddObserver(counter)

	g.Expect(observerError).To(BeNil())
	g.Expect(annealer.Observers()).To(ContainElement(counter))

	annealer.Anneal()

	g.Expect(counter.eventCounts[observer.StartedAnnealing]).To(BeNumerically("==", 1))
	g.Expect(counter.eventCounts[observer.FinishedAnnealing]).To(BeNumerically("==", 1))
	g.Expect(counter.eventCounts[observer.StartedIteration]).To(BeNumerically("==", expectedIterations))
	g.Expect(counter.eventCounts[observer.FinishedIteration]).To(BeNumerically("==", expectedIterations))
}

type CountingObserver struct {
	eventCounts map[observer.EventType]uint64
}

func (co *CountingObserver) ObserveEvent(event observer.Event) {
	co.eventCounts[event.EventType] += 1
}

func TestSimpleAnnealer_SetSolutionExplorer(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	const expectedTryCount = uint64(3)

	expectedParams := parameters.Map{
		MaximumIterations: int64(expectedTryCount),
	}

	annealer.SetParameters(expectedParams)

	expectedSolutionExplorer := new(TryCountingSolutionExplorer)

	explorerErr := annealer.SetSolutionExplorer(expectedSolutionExplorer)

	g.Expect(explorerErr).To(BeNil())
	g.Expect(annealer.SolutionExplorer()).To(BeIdenticalTo(expectedSolutionExplorer))

	annealer.Anneal()

	g.Expect(expectedSolutionExplorer.changesTried).To(BeNumerically("==", expectedTryCount))
}

type TryCountingSolutionExplorer struct {
	null.Explorer
	changesTried uint64
}

func (tcse *TryCountingSolutionExplorer) TryRandomChange() {
	tcse.changesTried += 1
}

func TestSimpleAnnealer_SetLogHandler(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	expectedParams := parameters.Map{
		MaximumIterations: int64(3),
	}

	annealer.SetParameters(expectedParams)

	expectedLogHandler := new(DummyLogHandler)

	annealer.SetLogHandler(expectedLogHandler)

	g.Expect(annealer.LogHandler()).To(BeIdenticalTo(expectedLogHandler))
}

func TestSimpleAnnealer_BadParameters(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	expectedParams := parameters.Map{
		MaximumIterations: "nope, not even close",
	}

	expectedLogHandler := new(DummyLogHandler)
	annealer.SetLogHandler(expectedLogHandler)

	actualError := annealer.SetParameters(expectedParams)

	g.Expect(actualError).To(Not(BeNil()))
	t.Log(actualError)
}

type DummyLogHandler struct {
	loggers.NullLogger
}

func TestSimpleAnnealer_PanicingExplorer(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	annealer.SetSolutionExplorer(new(flawedExplorer))

	expectedParams := parameters.Map{
		MaximumIterations: int64(1),
	}

	annealer.SetParameters(expectedParams)

	annealingCall := func() {
		annealer.Anneal()
	}

	g.Expect(annealingCall).To(Panic())
}

type flawedExplorer struct {
	null.Explorer
}

func (fe *flawedExplorer) TryRandomChange() {
	panic(errors.New("gotta panic"))
}
