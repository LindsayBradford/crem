// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"sort"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/pkg/errors"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
)

type SimpleAnnealer struct {
	explorer.ContainedExplorer
	model.ContainedModel

	loggers.ContainedLogger

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
	event := observer.NewEvent(observer.StartedAnnealing).
		WithId(sa.Id()).
		WithAttribute("MaximumIterations", sa.maximumIterations).
		WithAttribute("ObjectiveValue", sa.SolutionExplorer().ObjectiveValue()).
		WithAttribute("Temperature", sa.temperature).
		WithAttribute("CoolingFactor", sa.coolingFactor)

	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) iterationStarted() {
	sa.currentIteration++

	event := observer.NewEvent(observer.StartedIteration).
		WithId(sa.Id()).
		WithAttribute("CurrentIteration", sa.currentIteration).
		WithAttribute("MaximumIterations", sa.maximumIterations).
		WithAttribute("ObjectiveValue", sa.SolutionExplorer().ObjectiveValue()).
		WithAttribute("Temperature", sa.temperature)

	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) iterationFinished() {
	event := observer.NewEvent(observer.FinishedIteration).
		WithId(sa.Id()).
		WithAttribute("CurrentIteration", sa.currentIteration).
		WithAttribute("MaximumIterations", sa.maximumIterations).
		WithAttribute("ObjectiveValue", sa.SolutionExplorer().ObjectiveValue()).
		WithAttribute("ChangeInObjectiveValue", sa.SolutionExplorer().ChangeInObjectiveValue()).
		WithAttribute("ChangeIsDesirable", sa.SolutionExplorer().ChangeIsDesirable()).
		WithAttribute("AcceptanceProbability", sa.SolutionExplorer().AcceptanceProbability()).
		WithAttribute("ChangeAccepted", sa.SolutionExplorer().ChangeAccepted())

	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) annealingFinished() {
	sa.solution = *sa.fetchFinalModelSolution()

	event := observer.NewEvent(observer.FinishedAnnealing).
		WithId(sa.Id()).
		WithAttribute("CurrentIteration", sa.currentIteration).
		WithAttribute("MaximumIterations", sa.maximumIterations).
		WithAttribute("ObjectiveValue", sa.SolutionExplorer().ObjectiveValue()).
		WithAttribute("Temperature", sa.temperature)

	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) fetchFinalModelSolution() *solution.Solution {
	modelSolution := solution.NewSolution(sa.id)

	sa.addDecisionVariables(modelSolution)
	sa.addPlanningUnits(modelSolution)
	sa.addPlanningUnitManagementActionMap(modelSolution)

	return modelSolution
}

func (sa *SimpleAnnealer) addDecisionVariables(modelSolution *solution.Solution) {
	if sa.Model().DecisionVariables() == nil {
		return
	}

	solutionVariables := make(variable.EncodeableDecisionVariables, 0)

	for _, rawVariable := range *sa.Model().DecisionVariables() {
		solutionVariables = append(solutionVariables, variable.MakeEncodeable(rawVariable))
	}

	sort.Sort(solutionVariables)
	modelSolution.DecisionVariables = solutionVariables
}

func (sa *SimpleAnnealer) addPlanningUnits(modelSolution *solution.Solution) {
	if sa.Model().PlanningUnits() == nil {
		return
	}
	modelSolution.PlanningUnits = sa.Model().PlanningUnits()
}

func (sa *SimpleAnnealer) addPlanningUnitManagementActionMap(modelSolution *solution.Solution) {
	for _, action := range sa.Model().ActiveManagementActions() {
		planningUnit := solution.PlanningUnitId(action.PlanningUnit())
		actionType := solution.ManagementActionType(action.Type())
		modelSolution.ActiveManagementActions[planningUnit] = append(modelSolution.ActiveManagementActions[planningUnit], actionType)
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
