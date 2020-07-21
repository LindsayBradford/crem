package data

import "github.com/LindsayBradford/crem/internal/pkg/config/data"

type BasicScenarioConfig struct {
	Name string
}

type ScenarioConfig struct {
	Scenario BasicScenarioConfig
	Model    data.ModelConfig
}
