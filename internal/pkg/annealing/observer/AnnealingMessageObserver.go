// (c) 2018 Australian Rivers Institute.

package observer

import (
	"strconv"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/strings"
)

// AnnealingMessageObserver produces a stream of human-friendly, free-form text log entries from any observed
// Event instances received.
type AnnealingMessageObserver struct {
	AnnealingObserver
}

func (amo *AnnealingMessageObserver) WithLogHandler(handler logging.Logger) *AnnealingMessageObserver {
	amo.logHandler = handler
	return amo
}

func (amo *AnnealingMessageObserver) WithFilter(filter filters.Filter) *AnnealingMessageObserver {
	amo.filter = filter
	return amo
}

// ObserveEvent captures and converts Event instances into free-form text strings that it
// then passes onto its relevant Logger as an Info call.
func (amo *AnnealingMessageObserver) ObserveEvent(event observer.Event) {
	if amo.logHandler.BeingDiscarded(AnnealerLogLevel) || amo.filter.ShouldFilter(event) {
		return
	}

	var builder strings.FluentBuilder
	builder.
		Add("Id [", event.Id(), "],").
		Add("Event [", event.EventType.String(), "]: ")

	if observableAnnealer, isAnnealer := event.Source().(annealing.Observable); isAnnealer {
		amo.observeAnnealingEvent(observableAnnealer, event, &builder)
	} else {
		amo.observeEvent(event, &builder)
	}
}

func (amo *AnnealingMessageObserver) observeAnnealingEvent(observableAnnealer annealing.Observable, event observer.Event, builder *strings.FluentBuilder) {
	annealer := wrapAnnealer(observableAnnealer)
	explorer := wrapSolutionExplorer(observableAnnealer.ObservableExplorer())

	switch event.EventType {
	case observer.StartedAnnealing:
		builder.
			Add("Maximum Iterations [", annealer.MaximumIterations(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "], ").
			Add("Temperature [", annealer.Temperature(), "], ").
			Add("Cooling Factor [", annealer.CoolingFactor(), "]")
	case observer.StartedIteration:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaximumIterations(), "], ").
			Add("Temperature [", annealer.Temperature(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "]")
	case observer.FinishedIteration:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaximumIterations(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "], ").
			Add("Change [", explorer.ChangeInObjectiveValue(), "], ").
			Add("Desirable? [", explorer.ChangeIsDesirable(), "], ").
			Add("Acceptance Probability [", explorer.AcceptanceProbability(), "], ").
			Add("Accepted? [", explorer.ChangeAccepted(), "]")
	case observer.FinishedAnnealing:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaximumIterations(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "], ").
			Add("Temperature [", annealer.Temperature(), "]")
	default:
		// deliberately does nothing extra
	}

	amo.logHandler.LogAtLevel(AnnealerLogLevel, builder.String())
}

func (amo *AnnealingMessageObserver) observeEvent(event observer.Event, builder *strings.FluentBuilder) {
	switch event.EventType {
	case observer.Note:
		builder.Add("[", event.Note(), "]")
	case observer.ManagementAction:
		action, isAction := event.Source().(action.ManagementAction)
		if isAction {
			builder.
				Add("Type [", string(action.Type()), "], ").
				Add("Planning Unit [", action.PlanningUnit(), "], ").
				Add("Active [", strconv.FormatBool(action.IsActive()), "]")
		}
	default:
		// deliberately does nothing extra
	}

	amo.logHandler.LogAtLevel(model.LogLevel, builder.String())
}
