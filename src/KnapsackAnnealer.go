// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/logging"
)

var (
	Info    *log.Logger
)

func Init(infoHandle io.Writer) {
	Info = log.New(infoHandle, "",log.Ldate|log.Ltime|log.Lmicroseconds)
}

const ERROR_STATUS = 1

func main() {
	Init(os.Stdout)

	builder := new(AnnealerBuilder)

	annealer, err := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(1000).
		WithCoolingFactor(0.995).
		WithMaxIterations(5).
		WithObservers(NewNativeLogLibraryAnnealingLogger(Info)).
		Build()

	if (err != nil) {
		fmt.Println(err)
		os.Exit(ERROR_STATUS)
	}

	annealer.Anneal()
}

