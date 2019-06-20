// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	. "github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/model"
)

func NewInerpreter() *Interpreter {
	newInterpreter := new(Interpreter)

	newInterpreter.modelInterpreter = NewModelInterpreter()

	return newInterpreter
}

type Interpreter struct {
	modelInterpreter *ModelInterpreter
}

func (i *Interpreter) Interpret(config *Config) {
	i.modelInterpreter.Interpret(&config.Model)
}

func (i *Interpreter) Model() model.Model {
	return i.modelInterpreter.Model()
}
