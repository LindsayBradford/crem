// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	"log"

	. "github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/strings"
)

type NativeLogLibraryAnnealingLogger struct{
	logger *log.Logger
}

func NewNativeLogLibraryAnnealingLogger(logger *log.Logger) *NativeLogLibraryAnnealingLogger {
  newAnnealingLogger := &NativeLogLibraryAnnealingLogger {logger}
	return newAnnealingLogger
}

func (this *NativeLogLibraryAnnealingLogger) ObserveAnnealingEvent(event AnnealingEvent) {

	annealer := wrap(event.Annealer)

	var builder strings.FluentBuilder
	builder.Add("INFO {\"AnnealingEvent\": \"", event.EventType.String(), "\", ")

	switch event.EventType {
	case STARTED_ANNEALING:
		builder.
			Add("\"MaximumIterations\": ", annealer.MaxIterations(), ", ").
			Add("\"Temperature\": ", annealer.Temperature(), ", ").
			Add("\"CoolingFactor\": ", annealer.CoolingFactor())
	case STARTED_ITERATION, FINISHED_ANNEALING:
		builder.
			Add("\"CurrentIteration\": ", annealer.CurrentIteration(), ", ").
			Add("\"Temperature\": ", annealer.Temperature())
	case NOTE:
		builder.Add("\"Message\": \"", event.Note, "\"")
	default:
		// deliberately does nothing extra
	}

	builder.Add("}")

	this.logger.Println(builder.String())
}