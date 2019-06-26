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

const (
	NullModel               = "NullModel"
	DumbModel               = "DumbModel"
	MultiObjectiveDumbModel = "MultiObjectiveDumbModel"
	CatchmentModel          = "CatchmentModel"
)

type ModelConfigInterpreter struct {
	errors           *compositeErrors.CompositeError
	registeredModels map[string]ModelConfigFunction

	model model.Model
}

type ModelConfigFunction func(config ModelConfig) model.Model

func NewModelConfigInterpreter() *ModelConfigInterpreter {
	newInterpreter := new(ModelConfigInterpreter).initialise().
		RegisteringModel(
			NullModel,
			func(config ModelConfig) model.Model {
				return model.NewNullModel()
			},
		).RegisteringModel(
		DumbModel,
		func(config ModelConfig) model.Model {
			return dumb.NewModel().
				WithParameters(config.Parameters)
		},
	).RegisteringModel(
		MultiObjectiveDumbModel,
		func(config ModelConfig) model.Model {
			return modumb.NewModel().
				WithParameters(config.Parameters)
		},
	).RegisteringModel(
		CatchmentModel,
		func(config ModelConfig) model.Model {
			return catchment.NewModel().
				WithOleFunctionWrapper(threading.GetMainThreadChannel().Call).
				WithParameters(config.Parameters)
		},
	)

	return newInterpreter
}

func (i *ModelConfigInterpreter) initialise() *ModelConfigInterpreter {
	i.registeredModels = make(map[string]ModelConfigFunction, 0)
	i.errors = compositeErrors.New("Model Configuration")
	i.model = model.NullModel
	return i
}

func (i *ModelConfigInterpreter) Interpret(modelConfig *ModelConfig) *ModelConfigInterpreter {
	if _, foundModel := i.registeredModels[modelConfig.Type]; !foundModel {
		i.errors.Add(
			errors.New("configuration specifies a model type [\"" +
				modelConfig.Type + "\"], but no models are registered for that type"),
		)
		return i
	}

	configFunction := i.registeredModels[modelConfig.Type]
	newModel := configFunction(*modelConfig)
	if parameterisedModel, hasParameters := newModel.(parameters.Container); hasParameters {
		if paramErrors := parameterisedModel.ParameterErrors(); paramErrors != nil {
			wrappedErrors := errors.Wrap(paramErrors, "building model ["+modelConfig.Type+"]")
			i.errors.Add(wrappedErrors)
			return i
		}
	}
	i.model = newModel
	return i
}

func (i *ModelConfigInterpreter) RegisteringModel(modelType string, configFunction ModelConfigFunction) *ModelConfigInterpreter {
	i.registeredModels[modelType] = configFunction
	return i
}

func (i *ModelConfigInterpreter) Model() model.Model {
	return i.model
}

func (i *ModelConfigInterpreter) Errors() error {
	if i.errors.Size() > 0 {
		return i.errors
	}
	return nil
}
