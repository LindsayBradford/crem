// Copyright (c) 2019 Australian Rivers Institute.

package data

import "github.com/LindsayBradford/crem/internal/pkg/config/data"

type Config struct {
	MetaData data.MetaDataConfig

	Scenario ScenarioConfig
	Annealer data.AnnealerConfig
	Model    data.ModelConfig
}
