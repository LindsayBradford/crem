// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	. "github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/strings"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/modulators"
)

// FreeformAnnealingLogger produces a stream of human-friendly, free-form text log entries from any observed
// AnnealingEvent instances received.
type FreeformAnnealingLogger struct {
	AnnealingLogger
}

func (this *FreeformAnnealingLogger) WithLogHandler(handler LogHandler) *FreeformAnnealingLogger {
	this.logHandler = handler
	return this
}

func (this *FreeformAnnealingLogger) WithModulator(modulator  LoggingModulator) *FreeformAnnealingLogger {
	this.modulator = modulator
	return this
}

// ObserveAnnealingEvent captures and converts AnnealingEvent instances into free-form text strings that it
// then passes onto its relevant LogHandler as an Info call.
func (this *FreeformAnnealingLogger) ObserveAnnealingEvent(event AnnealingEvent) {
	if this.modulator.ShouldModulate(event) {
		return
	}

	annealer := wrap(event.Annealer)

	var builder strings.FluentBuilder
	builder.Add("Annealing Event [", event.EventType.String(), "]: ")

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
	case NOTE:
		builder.Add("[", event.Note, "]")
	default:
		// deliberately does nothing extra
	}

	this.logHandler.Info(builder.String())
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
