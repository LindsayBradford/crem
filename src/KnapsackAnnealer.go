// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"fmt"
	"os"

	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/logging"
)

const ERROR_STATUS = 1

func main() {
	builder := new(AnnealerBuilder)

	annealer, err := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(1000).
		WithCoolingFactor(0.995).
		WithMaxIterations(5).
		WithObservers(new(StdoutAnnealingLogger)).
		Build()

	if (err != nil) {
		fmt.Println(err)
		os.Exit(ERROR_STATUS)
	}

	annealer.Anneal()
}

