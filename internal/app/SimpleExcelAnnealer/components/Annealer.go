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

func BuildAnnealer(configuration *config.CRMConfig, humanLogHandler LogHandler, machineLogHandler LogHandler) Annealer {
	builder := new(AnnealerBuilder)
	machineAudienceObserver := new(AnnealingAttributeObserver).
		WithLogHandler(machineLogHandler).
		WithModulator(new(NullModulator))
		// WithModulator(new(IterationModuloLoggingModulator).WithModulo(200))
	humanAudienceObserver := new(AnnealingMessageObserver).
		WithLogHandler(humanLogHandler).
		// WithModulator(new(NullModulator))
		// WithModulator(new(IterationElapsedTimeLoggingModulator).WithWait(1 * time.Second))
		WithModulator(new(IterationModuloLoggingModulator).WithModulo(200))

	humanLogHandler.Debug("About to call AnnealerBuilder.Build() ")

	annealerConfig := configuration.Annealer

	newAnnealer, err := builder.
		AnnealerOfType(annealerConfig.Type).
		WithStartingTemperature(annealerConfig.StartingTemperature).
		WithCoolingFactor(annealerConfig.CoolingFactor).
		WithMaxIterations(annealerConfig.MaximumIterations).
		WithLogHandler(humanLogHandler).
		WithSolutionExplorer(new(SimpleExcelSolutionExplorer).WithPenalty(100)).
		WithEventNotifier(annealerConfig.EventNotifier).
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
