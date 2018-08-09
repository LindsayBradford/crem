// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"os"

	. "github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/config"
)

func BuildDumbAnnealer(annealerConfig *config.CRMConfig) Annealer {
	newAnnealer, logHandler, buildError :=
		new(config.AnnealerBuilder).
			WithConfig(annealerConfig).
			Build()

	if buildError != nil {
		logHandler.ErrorWithError(buildError)
		logHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	return newAnnealer
}
