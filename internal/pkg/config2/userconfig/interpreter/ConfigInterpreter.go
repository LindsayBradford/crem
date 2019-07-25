// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	. "github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
)

func NewInterpreter() *ConfigInterpreter {
	newInterpreter := new(ConfigInterpreter)

	newInterpreter.modelInterpreter = NewModelConfigInterpreter()
	newInterpreter.annealerInterpreter = NewAnnealerConfigInterpreter()
	newInterpreter.scenarioInterpreter = NewScenarioConfigInterpreter()

	return newInterpreter
}

type ConfigInterpreter struct {
	modelInterpreter *ModelConfigInterpreter
	model            model.Model

	annealerInterpreter *AnnealerConfigInterpreter
	annealer            annealing.Annealer

	scenarioInterpreter *ScenarioConfigInterpreter
	scenario            scenario.Scenario
}

func (i *ConfigInterpreter) Interpret(config *Config) *ConfigInterpreter {
	i.model = i.modelInterpreter.Interpret(&config.Model).Model()
	i.annealer = i.annealerInterpreter.Interpret(&config.Annealer).Annealer()
	i.scenario = i.scenarioInterpreter.Interpret(&config.Scenario).Scenario()

	i.annealer.SetModel(i.model)
	i.scenario.SetAnnealer(i.annealer)

	return i
}

func (i *ConfigInterpreter) Scenario() scenario.Scenario {
	assert.That(i.scenario != nil)
	return i.scenario
}
