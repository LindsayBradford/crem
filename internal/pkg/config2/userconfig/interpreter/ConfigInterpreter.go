// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/annealers"
	. "github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
)

func NewInterpreter() *ConfigInterpreter {
	newInterpreter := new(ConfigInterpreter).initialise()
	return newInterpreter
}

type ConfigInterpreter struct {
	errors *compositeErrors.CompositeError

	modelInterpreter *ModelConfigInterpreter
	model            model.Model

	annealerInterpreter *AnnealerConfigInterpreter
	annealer            annealing.Annealer

	scenarioInterpreter *ScenarioConfigInterpreter
	scenario            scenario.Scenario
}

func (i *ConfigInterpreter) initialise() *ConfigInterpreter {
	i.errors = compositeErrors.New("Configuration")

	i.model = model.NewNullModel()
	i.modelInterpreter = NewModelConfigInterpreter()

	i.annealer = &annealers.NullAnnealer{}
	i.annealerInterpreter = NewAnnealerConfigInterpreter()

	i.scenario = scenario.NullScenario
	i.scenarioInterpreter = NewScenarioConfigInterpreter()

	return i
}

func (i *ConfigInterpreter) Interpret(config *Config) *ConfigInterpreter {
	i.interpretModelConfig(&config.Model)
	i.interpretAnnealerConfig(&config.Annealer)
	i.interpretScenarioConfig(&config.Scenario)

	i.annealer.SetModel(i.model)
	i.scenario.SetAnnealer(i.annealer)

	return i
}

func (i *ConfigInterpreter) interpretModelConfig(config *ModelConfig) {
	i.model = i.modelInterpreter.Interpret(config).Model()
	if i.modelInterpreter.Errors() != nil {
		i.errors.Add(i.modelInterpreter.Errors())
	}
}

func (i *ConfigInterpreter) interpretAnnealerConfig(config *AnnealerConfig) {
	i.annealer = i.annealerInterpreter.Interpret(config).Annealer()
	if i.annealerInterpreter.Errors() != nil {
		i.errors.Add(i.annealerInterpreter.Errors())
	}
}

func (i *ConfigInterpreter) interpretScenarioConfig(config *ScenarioConfig) {
	i.scenario = i.scenarioInterpreter.Interpret(config).Scenario()
	if i.scenarioInterpreter.Errors() != nil {
		i.errors.Add(i.scenarioInterpreter.Errors())
	}
}

func (i *ConfigInterpreter) Scenario() scenario.Scenario {
	assert.That(i.scenario != nil)
	return i.scenario
}

func (i *ConfigInterpreter) Errors() error {
	if i.errors.Size() > 0 {
		return i.errors
	}
	return nil
}
