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

func BuildDumbAnnealer(logHandler LogHandler) Annealer {
	builder := new(AnnealerBuilder)
	humanAudienceObserver := new(AnnealingMessageObserver).
		WithLogHandler(logHandler).
		// WithModulator(new(NullModulator))
		WithModulator(new(IterationModuloLoggingModulator).WithModulo(100)) // No STARTED_ITERATION events, all FINISHED_ITERATION events

	logHandler.Debug("About to call AnnealerBuilder.Build() ")

	newAnnealer, err := builder.
		ElapsedTimeTrackingAnnealer().
		WithStartingTemperature(10).
		WithCoolingFactor(0.99).
		WithMaxIterations(1000000).
		WithDumbSolutionExplorer(100).
		WithLogHandler(logHandler).
		WithEventNotifier(new(SynchronousAnnealingEventNotifier)).
		// WithEventNotifier(new(ChanneledAnnealingEventNotifier)).
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
