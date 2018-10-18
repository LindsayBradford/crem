// (c) 2018 Australian Rivers Institute.

package logging

import (
	. "github.com/LindsayBradford/crm/annealing/shared"
	. "github.com/LindsayBradford/crm/logging/filters"
	. "github.com/LindsayBradford/crm/logging/handlers"
	"github.com/LindsayBradford/crm/strings"
)

// AnnealingMessageObserver produces a stream of human-friendly, free-form text log entries from any observed
// AnnealingEvent instances received.
type AnnealingMessageObserver struct {
	AnnealingLogger
}

func (amo *AnnealingMessageObserver) WithLogHandler(handler LogHandler) *AnnealingMessageObserver {
	amo.logHandler = handler
	return amo
}

func (amo *AnnealingMessageObserver) WithFilter(modulator LoggingFilter) *AnnealingMessageObserver {
	amo.filter = modulator
	return amo
}

// ObserveAnnealingEvent captures and converts AnnealingEvent instances into free-form text strings that it
// then passes onto its relevant LogHandler as an Info call.
func (amo *AnnealingMessageObserver) ObserveAnnealingEvent(event AnnealingEvent) {
	if amo.logHandler.BeingDiscarded(AnnealerLogLevel) || amo.filter.ShouldFilter(event) {
		return
	}

	annealer := wrapAnnealer(event.Annealer)
	explorer := wrapSolutionExplorer(event.Annealer.SolutionExplorer())

	var builder strings.FluentBuilder
	builder.Add("JobId [", event.Annealer.Id(), "], ", "Event [", event.EventType.String(), "]: ")

	switch event.EventType {
	case StartedAnnealing:
		builder.
			Add("Maximum Iterations [", annealer.MaxIterations(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "], ").
			Add("Temperature [", annealer.Temperature(), "], ").
			Add("Cooling Factor [", annealer.CoolingFactor(), "]")
	case StartedIteration:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaxIterations(), "], ").
			Add("Temperature [", annealer.Temperature(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "]")
	case FinishedIteration:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaxIterations(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "], ").
			Add("Change [", explorer.ChangeInObjectiveValue(), "], ").
			Add("Desirable? [", explorer.ChangeIsDesirable(), "], ").
			Add("Acceptance Probability [", explorer.AcceptanceProbability(), "], ").
			Add("Accepted? [", explorer.ChangeAccepted(), "]")
	case FinishedAnnealing:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaxIterations(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "], ").
			Add("Temperature [", annealer.Temperature(), "]")
	case Note:
		builder.Add("[", event.Note, "]")
	default:
		// deliberately does nothing extra
	}

	amo.logHandler.LogAtLevel(AnnealerLogLevel, builder.String())
}
