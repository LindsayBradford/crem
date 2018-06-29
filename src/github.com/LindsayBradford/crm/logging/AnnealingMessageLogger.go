// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	. "github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/strings"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/modulators"
)

// AnnealingMessageLogger produces a stream of human-friendly, free-form text log entries from any observed
// AnnealingEvent instances received.
type AnnealingMessageLogger struct {
	AnnealingLogger
}

func (this *AnnealingMessageLogger) WithLogHandler(handler LogHandler) *AnnealingMessageLogger {
	this.logHandler = handler
	return this
}

func (this *AnnealingMessageLogger) WithModulator(modulator  LoggingModulator) *AnnealingMessageLogger {
	this.modulator = modulator
	return this
}

// ObserveAnnealingEvent captures and converts AnnealingEvent instances into free-form text strings that it
// then passes onto its relevant LogHandler as an Info call.
func (this *AnnealingMessageLogger) ObserveAnnealingEvent(event AnnealingEvent) {
	if this.modulator.ShouldModulate(event) {
		return
	}

	annealer := wrap(event.Annealer)

	var builder strings.FluentBuilder
	builder.Add("Event [", event.EventType.String(), "]: ")

	switch event.EventType {
	case STARTED_ANNEALING:
		builder.
			Add("Maximum Iterations [", annealer.MaxIterations(), "], ").
			Add("Temperature [", annealer.Temperature(), "], ").
			Add("Cooling Factor [", annealer.CoolingFactor(), "]")
	case STARTED_ITERATION, FINISHED_ANNEALING:
		builder.
			Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaxIterations(), "], ").
			Add("Temperature [", annealer.Temperature(), "]")
	case NOTE, OBJECTIVE_EVALUATION:
		builder.Add("[", event.Note, "]")
	default:
		// deliberately does nothing extra
	}

	this.logHandler.LogAtLevel(ANNEALER, builder.String())
}

func wrap(eventAnnealer Annealer) *AnnealerStateFormatWrapper {
	wrapper := AnnealerStateFormatWrapper{
		AnnealerToFormat: eventAnnealer,
		MethodFormats: map[string]string{
			"Temperature":      "%0.4f",
			"CoolingFactor":    "%0.3f",
			"MaxIterations":    "%03d",
			"CurrentIteration": "%03d",
		},
	}
	return &wrapper
}
