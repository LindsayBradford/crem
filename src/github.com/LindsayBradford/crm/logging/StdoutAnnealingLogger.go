// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package logging

import (
	"fmt"
	. "github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/strings"
)

type StdoutAnnealingLogger struct{}

func (this *StdoutAnnealingLogger) ObserveAnnealingEvent(event AnnealingEvent) {

	var builder strings.FluentBuilder

	builder.Add("Annealing Event [", event.EventType.String(), "]: ")

	annealer := event.Annealer

	switch event.EventType {
	case STARTED_ANNEALING:
		builder.
			Add("Maximum Iterations [", uintToString(annealer.MaxIterations()), "], ").
			Add("Temperature [", float64ToString(annealer.Temperature()), "], ").
			Add("Cooling Factor [", float64ToString(annealer.CoolingFactor()), "]")
	case STARTED_ITERATION, FINISHED_ANNEALING:
		builder.
			Add("Iteration [", uintToString(annealer.CurrentIteration()), "/", uintToString(annealer.MaxIterations()), "], ").
			Add("Temperature [", float64ToString(annealer.Temperature()), "], ")
	case NOTE:
		builder.Add("[", event.Note, "]")
	default:
		// deliberately does nothing extra
	}

	fmt.Println(builder.String())
}

func uintToString(value uint) string {
	return fmt.Sprintf("%d", value)
}

func float64ToString(value float64) string {
	return fmt.Sprintf("%f", value)
}
