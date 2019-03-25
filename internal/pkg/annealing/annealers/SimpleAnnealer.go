// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/pkg/errors"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
)

type SimpleAnnealer struct {
	explorer.ContainedExplorer
	model.ModelContainer

	loggers.LoggerContainer

	parameters Parameters

	ContainedObservable
}

func (sa *SimpleAnnealer) Initialise() {
	sa.id = "Simple Annealer"
	sa.temperature = 1
	sa.currentIteration = 0
	sa.SetSolutionExplorer(null.NullExplorer)
	sa.SetLogHandler(new(loggers.NullLogger))
	sa.SetEventNotifier(new(observer.SynchronousAnnealingEventNotifier))
	sa.parameters.Initialise()
	sa.assignStateFromParameters()
}

func (sa *SimpleAnnealer) SetId(title string) {
	sa.ContainedObservable.SetId(title)
	sa.SolutionExplorer().SetId(title)
}

func (sa *SimpleAnnealer) SetParameters(params parameters.Map) error {
	sa.parameters.Merge(params)
	sa.assignStateFromParameters()
	return sa.parameters.ValidationErrors()
}

func (sa *SimpleAnnealer) assignStateFromParameters() {
	sa.SetTemperature(sa.parameters.GetFloat64(StartingTemperature))
	sa.coolingFactor = sa.parameters.GetFloat64(CoolingFactor)
	sa.maximumIterations = uint64(sa.parameters.GetInt64(MaximumIterations))
}

func (sa *SimpleAnnealer) ParameterErrors() error {
	return sa.parameters.ValidationErrors()
}

func (sa *SimpleAnnealer) DeepClone() annealing.Annealer {
	clone := *sa
	explorerClone := sa.SolutionExplorer().DeepClone()
	clone.SetSolutionExplorer(explorerClone)
	return &clone
}

func (sa *SimpleAnnealer) SetTemperature(temperature float64) error {
	if temperature <= 0 {
		return errors.New("invalid attempt to set annealer temperature to value <= 0")
	}
	sa.temperature = temperature
	return nil
}

func (sa *SimpleAnnealer) Model() model.Model {
	return sa.SolutionExplorer().Model()
}

func (sa *SimpleAnnealer) notifyObservers(eventType observer.EventType) {
	eventSource := sa.CloneObservable()
	event := observer.NewEvent(eventType).
		WithId(sa.Id()).
		WithSource(eventSource)
	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) CloneObservable() annealing.Observable {
	observable := sa.ContainedObservable
	explorerClone := sa.SolutionExplorer().CloneObservable()
	observable.SetObservableExplorer(explorerClone)
	return &observable
}

func (sa *SimpleAnnealer) Anneal() {
	defer sa.handlePanicRecovery()

	sa.SolutionExplorer().Initialise()
	defer sa.SolutionExplorer().TearDown()

	sa.annealingStarted()

	for done := sa.initialDoneValue(); !done; {
		sa.iterationStarted()

		sa.SolutionExplorer().TryRandomChange(sa.temperature)

		sa.iterationFinished()
		sa.cooldown()
		done = sa.checkIfDone()
	}

	sa.annealingFinished()
}

func (sa *SimpleAnnealer) handlePanicRecovery() {
	if r := recover(); r != nil {
		if rAsError, isError := r.(error); isError {
			wrappingError := errors.Wrap(rAsError, "annealing function failed")
			panic(wrappingError)
		}
		panic(r)
	}
}

func (sa *SimpleAnnealer) annealingStarted() {
	sa.notifyObservers(observer.StartedAnnealing)
}

func (sa *SimpleAnnealer) iterationStarted() {
	sa.currentIteration++
	sa.notifyObservers(observer.StartedIteration)
}

func (sa *SimpleAnnealer) iterationFinished() {
	sa.notifyObservers(observer.FinishedIteration)
}

func (sa *SimpleAnnealer) annealingFinished() {
	sa.solution = *sa.fetchFinalModelSolution()
	sa.notifyObservers(observer.FinishedAnnealing)
}

func (sa *SimpleAnnealer) fetchFinalModelSolution() *solution.Solution {
	modelSolution := new(solution.Solution)
	modelSolution.Id = sa.id

	sa.addDecisionVariables(modelSolution)

	return modelSolution
}

func (sa *SimpleAnnealer) addDecisionVariables(modelSolution *solution.Solution) {
	modelSolution.DecisionVariables = make(attributes.Attributes, 0)

	if sa.Model().DecisionVariables() == nil {
		return
	}

	for _, variable := range *sa.Model().DecisionVariables() {
		newPair := attributes.NameValuePair{
			Name:  variable.Name(),
			Value: variable.Value(),
		}
		modelSolution.DecisionVariables = append(modelSolution.DecisionVariables, newPair)
	}
}

func (sa *SimpleAnnealer) initialDoneValue() bool {
	return sa.parameters.GetInt64(MaximumIterations) == 0
}

func (sa *SimpleAnnealer) checkIfDone() bool {
	return sa.currentIteration >= uint64(sa.parameters.GetInt64(MaximumIterations))
}

func (sa *SimpleAnnealer) cooldown() {
	sa.temperature *= sa.parameters.GetFloat64(CoolingFactor)
}
