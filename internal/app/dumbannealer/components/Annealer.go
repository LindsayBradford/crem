// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"

	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/annealing/logging"
	. "github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/config"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/modulators"
)

func BuildDumbAnnealer(config *config.CRMConfig, logHandler LogHandler) Annealer {
	builder := new(AnnealerBuilder)
	humanAudienceObserver := new(AnnealingMessageObserver).
		WithLogHandler(logHandler).
		// WithModulator(new(NullModulator))
		WithModulator(new(IterationModuloLoggingModulator).WithModulo(100)) // No StartedIteration events, all FinishedIteration events

	logHandler.Debug("About to call AnnealerBuilder.Build() ")

	newAnnealer, err := builder.
		ElapsedTimeTrackingAnnealer().
		WithStartingTemperature(config.Annealer.StartingTemperature).
		WithCoolingFactor(config.Annealer.CoolingFactor).
		WithMaxIterations(config.Annealer.MaximumIterations).
		WithDumbSolutionExplorer(100).
		WithLogHandler(logHandler).
		WithEventNotifier(config.Annealer.EventNotifier).
		WithObservers(humanAudienceObserver).
		Build()

	logHandler.Debug("Call to AnnealerBuilder.Build() finished")

	if err != nil {
		logHandler.ErrorWithError(err)
		logHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	return newAnnealer
}
