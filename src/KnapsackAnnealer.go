// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package main

import "Annealer"

func main() {
	builder := new(Annealer.AnnealerBuilder)

	annealer := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(1000).
		WithCoolingFactor(0.995).
		WithMaxIterations(5).
		WithObservers(new(Annealer.StdoutAnnealingObserver)).
		Build()

	annealer.Anneal()
}
