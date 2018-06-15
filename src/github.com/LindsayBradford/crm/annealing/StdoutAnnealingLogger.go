// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import "fmt"

type StdoutAnnealingLogger struct{}

func (this *StdoutAnnealingLogger) ObserveAnnealingEvent(event AnnealingEvent) {
	fmt.Printf("Annealing Event [%s]: ", event.eventType.String())

	annealer := event.annealer

	switch event.eventType {
	case STARTED_ANNEALING:
		fmt.Printf("Maximum Iterations [%d], Temperature [%f], Cooling Factor [%f], ",
			annealer.MaxIterations(), annealer.Temperature(), annealer.CoolingFactor())
	case STARTED_ITERATION, FINISHED_ANNEALING:
		fmt.Printf("Iteration [%d/%d], Temperature [%f]",
			annealer.CurrentIteration(), annealer.MaxIterations(), annealer.Temperature())
	case NOTE:
		fmt.Printf("[%s]", event.note)
	default:
		// deliberately does nothing extra
	}

	fmt.Println()
}
