// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"testing"

	"github.com/LindsayBradford/crm/annealing/objectives"
	"github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/onsi/gomega"
)

func TestSimpleAnnealer_Initialise(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	g.Expect(
		annealer.Temperature()).To(BeIdenticalTo(float64(1)),
		"Annealer should have built with default Temperature of 1")

	g.Expect(
		annealer.CoolingFactor()).To(BeIdenticalTo(float64(1)),
		"Annealer should have built with default Cooling Factor of 1")

	g.Expect(
		annealer.MaxIterations()).To(BeZero(),
		"Annealer should have built with default iterations of 0")

	g.Expect(
		annealer.CurrentIteration()).To(BeZero(),
		"Annealer should have built with current iteration of 0")

	g.Expect(
		annealer.LogHandler()).To(Equal(handlers.NULL_LOG_HANDLER),
		"Annealer should have built with nullLogHandler")

	g.Expect(
		annealer.ObjectiveManager()).To(Equal(objectives.NULL_OBJECTIVE_MANAGER),
		"Annealer should have built with nullObjectiveManager")

	g.Expect(
		annealer.Observers()).To(BeNil(),
		"Annealer should have built with no AnnealerObservers")
}

func TestSimpleAnnealer_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	tempErr := annealer.SetTemperature(-1)

	g.Expect(tempErr).To(Not(BeNil()))
	g.Expect(
		annealer.Temperature()).To(BeIdenticalTo(float64(1)),
		"Annealer should have ignored crap Temperature set attempt")

	coolingFactorErr := annealer.SetCoolingFactor(1.000001)

	g.Expect(coolingFactorErr).To(Not(BeNil()))
	g.Expect(annealer.CoolingFactor()).To(BeIdenticalTo(float64(1)),
		"Annealer should have ignored crap CoolingFactor set attempt")

	logHandlerErr := annealer.SetLogHandler(nil)

	g.Expect(logHandlerErr).To(Not(BeNil()))
	g.Expect(annealer.LogHandler()).To(Equal(handlers.NULL_LOG_HANDLER),
		"Annealer should have ignored crap LogHandler set attempt")

	objectiveMgrErr := annealer.SetObjectiveManager(nil)

	g.Expect(objectiveMgrErr).To(Not(BeNil()))
	g.Expect(annealer.ObjectiveManager()).To(Equal(objectives.NULL_OBJECTIVE_MANAGER),
		"Annealer should have ignored crap Objective Manager set attempt")

	observersErr := annealer.AddObserver(nil)

	g.Expect(observersErr).To(Not(BeNil()))
	g.Expect(annealer.Observers()).To(BeNil(),
		"Annealer should have ignored crap AnnealerObserver set attempt")
}


func TestSimpleAnnealer_Anneal(t *testing.T) {
	g := NewGomegaWithT(t)

	const startTemperature float64 = 1000.0
	const coolingFactor float64 = 0.5
	const iterations uint = 3
	const expectedEndTemperature float64 = ((startTemperature * coolingFactor) * coolingFactor) * coolingFactor

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	annealer.SetTemperature(startTemperature)
	annealer.SetCoolingFactor(coolingFactor)
	annealer.SetMaxIterations(iterations)

	g.Expect(
		annealer.CurrentIteration()).To(BeZero(),
		"Annealer should have started with current iteration of 0")

	g.Expect(
		annealer.Temperature()).To(BeIdenticalTo(startTemperature),
		"Annealer should have started with expected start temperature")

	annealer.Anneal()

	g.Expect(
		annealer.CurrentIteration()).To(BeIdenticalTo(annealer.MaxIterations()),
		"Annealer should have ended with current iteration = max iterations")

	g.Expect(
		annealer.Temperature()).To(BeIdenticalTo(expectedEndTemperature),
		"Annealer should have ended with temperature modified by cooling factor * iterations")
}

type CountingObserver struct {
	eventCounts map[AnnealingEventType] uint
}

func (this *CountingObserver) ObserveAnnealingEvent(event AnnealingEvent) {
	this.eventCounts[event.EventType] += 1
}

func TestSimpleAnnealer_AddObserver(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	const expectedIterations = uint(3)

	annealer.SetTemperature(1000.0)
	annealer.SetCoolingFactor(0.5)
	annealer.SetMaxIterations(expectedIterations)

	countingObserver := new(CountingObserver)
	countingObserver.eventCounts = make(map[AnnealingEventType] uint)

	observerError := annealer.AddObserver(countingObserver)

	g.Expect(observerError).To(BeNil())
	g.Expect(annealer.Observers()).To(ContainElement(countingObserver),
		"Annealer should have accepted CountingObserver as new AnnealerObserver")

	annealer.Anneal()

	g.Expect(countingObserver.eventCounts[STARTED_ANNEALING]).To(BeIdenticalTo(uint(1)),
		"Annealer should have posted 1 STARTED_ANNEALING event")

	g.Expect(countingObserver.eventCounts[FINISHED_ANNEALING]).To(BeIdenticalTo(uint(1)),
		"Annealer should have posted 1 FINISHED_ANNEALING event")

	g.Expect(countingObserver.eventCounts[STARTED_ITERATION]).To(BeIdenticalTo(expectedIterations),
		"Annealer should have posted <expectedIterations> of  STARTED_ITERATION event")

	g.Expect(countingObserver.eventCounts[FINISHED_ITERATION]).To(BeIdenticalTo(expectedIterations),
		"Annealer should have posted <expectedIterations> of  FINISHED_ITERATION event")
}

type TryCountingObjectiveManager struct {
	objectives.NullObjectiveManager
	changesTried uint
}

func (this *TryCountingObjectiveManager) TryRandomChange(temperature float64) {
	this.changesTried += 1
}

func TestSimpleAnnealer_SetObjectiveManager(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	const expectedObjectiveManagerTryCount = uint(3)

	annealer.SetTemperature(1000.0)
	annealer.SetCoolingFactor(0.5)
	annealer.SetMaxIterations(expectedObjectiveManagerTryCount)

	expectedObjectiveManager := new(TryCountingObjectiveManager)

	objectiveManagerErr := annealer.SetObjectiveManager(expectedObjectiveManager)

	g.Expect(objectiveManagerErr).To(BeNil())
	g.Expect(annealer.objectiveManager).To(BeIdenticalTo(expectedObjectiveManager),
		"Annealer should have accepted CountingObserver as new ObjectiveManager")

	annealer.Anneal()

	g.Expect(expectedObjectiveManager.changesTried).To(BeIdenticalTo(expectedObjectiveManagerTryCount),
		"Annealer should have tried same number of changes as iterations")
}

type DummyLogHandler struct {
	 handlers.NullLogHandler
}

func TestSimpleAnnealer_SetLogHandler(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	const expectedObjectiveManagerTryCount = uint(3)

	annealer.SetTemperature(1000.0)
	annealer.SetCoolingFactor(0.5)
	annealer.SetMaxIterations(expectedObjectiveManagerTryCount)

	expectedLogHandler := new(DummyLogHandler)

	logHandlerErr := annealer.SetLogHandler(expectedLogHandler)

	g.Expect(logHandlerErr).To(BeNil())
	g.Expect(annealer.logger).To(BeIdenticalTo(expectedLogHandler),
		"Annealer should have accepted DummyLogHandler as new logger")
}