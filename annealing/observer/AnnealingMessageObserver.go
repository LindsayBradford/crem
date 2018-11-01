// (c) 2018 Australian Rivers Institute.

package observer

import (
	"github.com/LindsayBradford/crem/annealing"
	"github.com/LindsayBradford/crem/annealing/observer/filters"
	"github.com/LindsayBradford/crem/logging"
	"github.com/LindsayBradford/crem/strings"
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

// ObserveAnnealingEvent captures and converts Event instances into free-form text strings that it
// then passes onto its relevant Logger as an Info call.
func (amo *AnnealingMessageObserver) ObserveAnnealingEvent(event annealing.Event) {
	if amo.logHandler.BeingDiscarded(AnnealerLogLevel) || amo.filter.ShouldFilter(event) {
		return
	}

	annealer := wrapAnnealer(event.Annealer)
	explorer := wrapSolutionExplorer(event.Annealer.SolutionExplorer())

	var builder strings.FluentBuilder
	builder.Add("Id [", event.Annealer.Id(), "], ", "Event [", event.EventType.String(), "]: ")

	switch event.EventType {
	case annealing.StartedAnnealing:
		builder.
			Add("Maximum Iterations [", annealer.MaxIterations(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "], ").
			Add("Temperature [", annealer.Temperature(), "], ").
			Add("Cooling Factor [", annealer.CoolingFactor(), "]")
	case annealing.StartedIteration:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaxIterations(), "], ").
			Add("Temperature [", annealer.Temperature(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "]")
	case annealing.FinishedIteration:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaxIterations(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "], ").
			Add("Change [", explorer.ChangeInObjectiveValue(), "], ").
			Add("Desirable? [", explorer.ChangeIsDesirable(), "], ").
			Add("Acceptance Probability [", explorer.AcceptanceProbability(), "], ").
			Add("Accepted? [", explorer.ChangeAccepted(), "]")
	case annealing.FinishedAnnealing:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaxIterations(), "], ").
			Add("Objective value [", explorer.ObjectiveValue(), "], ").
			Add("Temperature [", annealer.Temperature(), "]")
	case annealing.Note:
		builder.Add("[", event.Note, "]")
	default:
		// deliberately does nothing extra
	}

	amo.logHandler.LogAtLevel(AnnealerLogLevel, builder.String())
}
