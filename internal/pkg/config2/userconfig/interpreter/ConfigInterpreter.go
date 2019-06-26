// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	. "github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/model"
)

func NewInerpreter() *ConfigInterpreter {
	newInterpreter := new(ConfigInterpreter)

	newInterpreter.modelInterpreter = NewModelConfigInterpreter()
	newInterpreter.annealerInterpreter = NewAnnealerConfigInterpreter()

	return newInterpreter
}

type ConfigInterpreter struct {
	modelInterpreter    *ModelConfigInterpreter
	annealerInterpreter *AnnealerConfigInterpreter
}

func (i *ConfigInterpreter) Interpret(config *Config) {
	i.modelInterpreter.Interpret(&config.Model)
	i.annealerInterpreter.Interpret(&config.Annealer)

	i.Annealer().SetModel(i.Model())
}

func (i *ConfigInterpreter) Model() model.Model {
	return i.modelInterpreter.Model()
}

func (i *ConfigInterpreter) Annealer() annealing.Annealer {
	return i.annealerInterpreter.Annealer()
}
