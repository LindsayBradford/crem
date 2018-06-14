// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package Annealer

import "fmt"

type AnnealingObserver interface {
	ObserveAnnealingEvent(event AnnealingEvent, annealer Annealer)
}

type StdoutAnnealingObserver struct{}

func (this *StdoutAnnealingObserver) ObserveAnnealingEvent(event AnnealingEvent, annealer Annealer) {
	fmt.Printf("Event [%s] received.", event.String())

	switch event {
	case STARTED_ANNEALING:
		fmt.Printf(" Temperature = [%f], MaxIterations = [%d]", annealer.Temperature(), annealer.MaxIterations())
	case STARTED_ITERATION, FINISHED_ANNEALING:
		fmt.Printf(" Temperature = [%f], CurrIteration = [%d/%d]", annealer.Temperature(), annealer.CurrentIteration(), annealer.MaxIterations())
	default:
		// deliberately does nothing extra
	}

	fmt.Println()
}
