// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"

	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/annealing/logging"
	. "github.com/LindsayBradford/crm/annealing/shared"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/modulators"
)

func BuildAnnealer(humanLogHandler LogHandler, machineLogHandler LogHandler) Annealer {
	builder := new(AnnealerBuilder)
	machineAudienceObserver := new(AnnealingAttributeObserver).
		WithLogHandler(machineLogHandler).
		WithModulator(new(NullModulator))
	humanAudienceObserver := new(AnnealingMessageObserver).
		WithLogHandler(humanLogHandler).
	// WithModulator(new(NullModulator))
	// WithModulator(new(IterationElapsedTimeLoggingModulator).WithWait(1 * time.Second))
		WithModulator(new(IterationModuloLoggingModulator).WithModulo(200))

	humanLogHandler.Debug("About to call AnnealerBuilder.Build() ")

	newAnnealer, err := builder.
		ElapsedTimeTrackingAnnealer().
		WithLogHandler(humanLogHandler).
		WithObjectiveManager(new(KnapsackObjectiveManager)).
		WithStartingTemperature(10).
		WithCoolingFactor(0.997).
		WithMaxIterations(2000).
		WithObservers(machineAudienceObserver, humanAudienceObserver).
		Build()

	humanLogHandler.Debug("Call to AnnealerBuilder.Build() finished")

	if err != nil {
		humanLogHandler.ErrorWithError(err)
		humanLogHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	return newAnnealer
}

