// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	. "github.com/LindsayBradford/crm/annealing/shared"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/modulators"
	"github.com/LindsayBradford/crm/strings"
)

// AnnealingMessageObserver produces a stream of human-friendly, free-form text log entries from any observed
// AnnealingEvent instances received.
type AnnealingMessageObserver struct {
	AnnealingLogger
}

func (this *AnnealingMessageObserver) WithLogHandler(handler LogHandler) *AnnealingMessageObserver {
	this.logHandler = handler
	return this
}

func (this *AnnealingMessageObserver) WithModulator(modulator LoggingModulator) *AnnealingMessageObserver {
	this.modulator = modulator
	return this
}

// ObserveAnnealingEvent captures and converts AnnealingEvent instances into free-form text strings that it
// then passes onto its relevant LogHandler as an Info call.
func (this *AnnealingMessageObserver) ObserveAnnealingEvent(event AnnealingEvent) {
	if this.logHandler.BeingDiscarded(ANNEALER) || this.modulator.ShouldModulate(event) {
		return
	}

	annealer := wrapAnnealer(event.Annealer)
	objectiveManager := wrapObjectiveManager(event.Annealer.ObjectiveManager())

	var builder strings.FluentBuilder
	builder.Add("Event [", event.EventType.String(), "]: ")

	switch event.EventType {
	case STARTED_ANNEALING:
		builder.
			Add("Maximum Iterations [", annealer.MaxIterations(), "], ").
			Add("Objective value [", objectiveManager.ObjectiveValue(), "], ").
			Add("Temperature [", annealer.Temperature(), "], ").
			Add("Cooling Factor [", annealer.CoolingFactor(), "]")
	case STARTED_ITERATION:
		objectiveManager := wrapObjectiveManager(event.Annealer.ObjectiveManager())
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaxIterations(), "], ").
			Add("Temperature [", annealer.Temperature(), "], ").
			Add("Objective value [", objectiveManager.ObjectiveValue(), "]")
	case FINISHED_ITERATION:
		objectiveManager := wrapObjectiveManager(event.Annealer.ObjectiveManager())
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaxIterations(), "], ").
			Add("Objective value [", objectiveManager.ObjectiveValue(), "], ").
			Add("Change [", objectiveManager.ChangeInObjectiveValue(), "], ").
			Add("Desirable? [", objectiveManager.ChangeIsDesirable(), "], ").
			Add("Acceptance Probability [", objectiveManager.AcceptanceProbability(), "], ").
			Add("Accepted? [", objectiveManager.ChangeAccepted(), "]")
	case FINISHED_ANNEALING:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaxIterations(), "], ").
			Add("Objective value [", objectiveManager.ObjectiveValue(), "], ").
			Add("Temperature [", annealer.Temperature(), "]")
	case NOTE:
		builder.Add("[", event.Note, "]")
	default:
		// deliberately does nothing extra
	}

	this.logHandler.LogAtLevel(ANNEALER, builder.String())
}