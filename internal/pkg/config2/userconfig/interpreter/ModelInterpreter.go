// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	. "github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/dumb"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

type ModelInterpreter struct {
	errors           *compositeErrors.CompositeError
	registeredModels map[string]ModelConfigFunction

	model model.Model
}

type ModelConfigFunction func(config ModelConfig) model.Model

func NewModelInterpreter() *ModelInterpreter {
	newInterpreter := new(ModelInterpreter).
		RegisteringModel(
			"NullModel",
			func(config ModelConfig) model.Model {
				return model.NewNullModel()
			},
		).RegisteringModel(
		"DumbModel",
		func(config ModelConfig) model.Model {
			return dumb.NewModel().WithParameters(config.Parameters)
		},
	).RegisteringModel(
		"MultiObjectiveDumbModel",
		func(config ModelConfig) model.Model {
			return modumb.NewModel().
				WithParameters(config.Parameters)
		},
	).RegisteringModel(
		"CatchmentModel",
		func(config ModelConfig) model.Model {
			return catchment.NewModel().
				WithOleFunctionWrapper(threading.GetMainThreadChannel().Call).
				WithParameters(config.Parameters)
		},
	)

	newInterpreter.errors = compositeErrors.New("Model Configuration")

	return newInterpreter
}

func (i *ModelInterpreter) Interpret(modelConfig *ModelConfig) {
	if _, foundModel := i.registeredModels[modelConfig.Type]; !foundModel {
		i.errors.Add(
			errors.New("configuration specifies a model type [\"" +
				modelConfig.Type + "\"], but no models are registered for that type"),
		)
		return
	}

	configFunction := i.registeredModels[modelConfig.Type]
	newModel := configFunction(*modelConfig)
	if parameterisedModel, hasParameters := newModel.(parameters.Container); hasParameters {
		if paramErrors := parameterisedModel.ParameterErrors(); paramErrors != nil {
			wrappedErrors := errors.Wrap(paramErrors, "building model ["+modelConfig.Type+"]")
			i.errors.Add(wrappedErrors)
		}
	}
}

func (i *ModelInterpreter) RegisteringModel(modelType string, configFunction ModelConfigFunction) *ModelInterpreter {
	i.registeredModels[modelType] = configFunction
	return i
}

func (i *ModelInterpreter) Model() model.Model {
	return i.model
}

func (i *ModelInterpreter) Errors() error {
	return i.errors
}
