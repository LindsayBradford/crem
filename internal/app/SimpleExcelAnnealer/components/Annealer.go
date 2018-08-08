// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"

	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/config"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

func BuildAnnealer(configuration *config.CRMConfig, humanLogHandler LogHandler, explorer solution.SolutionExplorer, observers ...AnnealingObserver) Annealer {
	builder := new(AnnealerBuilder)

	humanLogHandler.Debug("About to call AnnealerBuilder.Build() ")

	annealerConfig := configuration.Annealer

	newAnnealer, err := builder.
		AnnealerOfType(annealerConfig.Type).
		WithStartingTemperature(annealerConfig.StartingTemperature).
		WithCoolingFactor(annealerConfig.CoolingFactor).
		WithMaxIterations(annealerConfig.MaximumIterations).
		WithLogHandler(humanLogHandler).
		WithSolutionExplorer(explorer).
		WithEventNotifier(annealerConfig.EventNotifier).
		WithObservers(observers...).
		Build()

	humanLogHandler.Debug("Call to AnnealerBuilder.Build() finished")

	if err != nil {
		humanLogHandler.ErrorWithError(err)
		humanLogHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	return newAnnealer
}
