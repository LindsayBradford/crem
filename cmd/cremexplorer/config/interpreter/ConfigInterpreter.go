// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	data2 "github.com/LindsayBradford/crem/cmd/cremexplorer/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/annealers"
	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/config/interpreter"
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

	modelInterpreter *interpreter.ModelConfigInterpreter
	model            model.Model

	annealerInterpreter *interpreter.AnnealerConfigInterpreter
	annealer            annealing.Annealer

	scenarioInterpreter *ScenarioConfigInterpreter
	scenario            scenario.Scenario
}

func (i *ConfigInterpreter) initialise() *ConfigInterpreter {
	i.errors = compositeErrors.New("Configuration")

	i.model = model.NewNullModel()
	i.modelInterpreter = interpreter.NewModelConfigInterpreter()

	i.annealer = &annealers.NullAnnealer{}
	i.annealerInterpreter = interpreter.NewAnnealerConfigInterpreter()

	i.scenario = scenario.NullScenario
	i.scenarioInterpreter = NewScenarioConfigInterpreter()

	return i
}

func (i *ConfigInterpreter) Interpret(config *data2.Config) *ConfigInterpreter {
	i.interpretModelConfig(&config.Model)
	i.interpretAnnealerConfig(&config.Annealer)
	i.interpretScenarioConfig(&config.Scenario)

	i.annealer.SetModel(i.model)
	i.scenario.SetAnnealer(i.annealer)

	return i
}

func (i *ConfigInterpreter) interpretModelConfig(config *data.ModelConfig) {
	i.model = i.modelInterpreter.Interpret(config).Model()
	if i.modelInterpreter.Errors() != nil {
		i.errors.Add(i.modelInterpreter.Errors())
	}
}

func (i *ConfigInterpreter) interpretAnnealerConfig(config *data.AnnealerConfig) {
	i.annealer = i.annealerInterpreter.Interpret(config).Annealer()
	if i.annealerInterpreter.Errors() != nil {
		i.errors.Add(i.annealerInterpreter.Errors())
	}
}

func (i *ConfigInterpreter) interpretScenarioConfig(config *data2.ScenarioConfig) {
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
