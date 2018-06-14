// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package Annealer

import "fmt"

type AnnealingObserver interface {
	ObserveAnnealingEvent(event AnnealingEvent)
}

type StdoutAnnealingObserver struct{}

func (this *StdoutAnnealingObserver) ObserveAnnealingEvent(event AnnealingEvent) {
	fmt.Printf("Event [%s] received.", event.eventType.String())

	annealer := event.annealer

	switch event.eventType {
	case STARTED_ANNEALING:
		fmt.Printf(" Temperature = [%f], Maximum Iterations = [%d]", annealer.Temperature(), annealer.MaxIterations())
	case STARTED_ITERATION, FINISHED_ANNEALING:
		fmt.Printf(" Temperature = [%f], Current Iteration = [%d/%d]", annealer.Temperature(), annealer.CurrentIteration(), annealer.MaxIterations())
	case NOTE:
		fmt.Printf(" [%s]", event.note)
	default:
		// deliberately does nothing extra
	}

	fmt.Println()
}
