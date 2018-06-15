// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import . "github.com/LindsayBradford/crm/annealing"
import . "github.com/LindsayBradford/crm/logging"

func main() {
	builder := new(AnnealerBuilder)

	annealer := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(1000).
		WithCoolingFactor(0.995).
		WithMaxIterations(5).
		WithObservers(new(StdoutAnnealingLogger)).
		Build()

	annealer.Anneal()
}
