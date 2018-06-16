// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package logging

import (
	"fmt"
	. "github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/strings"
)

type StdoutAnnealingLogger struct {}

func (this *StdoutAnnealingLogger) ObserveAnnealingEvent(event AnnealingEvent) {
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
			Add("Temperature [", annealer.Temperature(), "], ")
	case NOTE:
		builder.Add("[", event.Note, "]")
	default:
		// deliberately does nothing extra
	}

	fmt.Println(builder.String())
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
