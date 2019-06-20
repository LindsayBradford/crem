// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"errors"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/dumb"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	. "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/threading"
	errors2 "github.com/pkg/errors"
)

type modelBuilder struct {
	errors           *CompositeError
	config           []ModelConfig
	registeredModels map[string]ModelConfigFunction
}

type ModelConfigFunction func(config ModelConfig) model.Model

func (builder *modelBuilder) initialise() *modelBuilder {
	if builder.errors == nil {
		builder.errors = New("Model initialisation")
	}

	if builder.registeredModels == nil {
		builder.registerBaseModels()
	}

	return builder
}

func (builder *modelBuilder) registerBaseModels() {
	builder.registeredModels = make(map[string]ModelConfigFunction, 2)

	builder.RegisteringModel(
		"NullModel",
		func(config ModelConfig) model.Model {
			return model.NewNullModel().WithName(config.Name)
		},
	)

	builder.RegisteringModel(
		"DumbModel",
		func(config ModelConfig) model.Model {
			return dumb.NewModel().WithName(config.Name).WithParameters(config.Parameters)
		},
	)

	builder.RegisteringModel(
		"CatchmentModel",
		func(config ModelConfig) model.Model {
			return catchment.NewModel().
				WithName(config.Name).
				WithOleFunctionWrapper(threading.GetMainThreadChannel().Call).
				WithParameters(config.Parameters)
		},
	)
}

func (builder *modelBuilder) WithConfig(cremConfig *CREMConfig) *modelBuilder {
	builder.initialise()
	builder.config = cremConfig.Models
	return builder
}

func (builder *modelBuilder) RegisteringModel(modelType string, configFunction ModelConfigFunction) *modelBuilder {
	builder.initialise()
	builder.registeredModels[modelType] = configFunction
	return builder
}

func (builder *modelBuilder) Build(modelName string) (model.Model, error) {
	var myModel model.Model
	if len(builder.config) == 0 {
		builder.errors.Add(errors.New("configuration failed to specify any models"))
	} else {
		myModel = builder.findMyModel(modelName, builder.buildModels())
	}

	if builder.errors.Size() > 0 {
		return nil, builder.errors
	}

	return myModel, nil
}

func (builder *modelBuilder) findMyModel(myModelName string, models []model.Model) model.Model {
	for _, model := range models {
		if model != nil && model.Name() == myModelName {
			return model
		}
	}
	builder.errors.Add(
		errors.New("configuration specifies a non-existent model [\"" +
			myModelName + "\"] for its Annealer"),
	)
	return nil
}

func (builder *modelBuilder) buildModels() []model.Model {
	modelList := make([]model.Model, len(builder.config))
	for index, currConfig := range builder.config {
		_, foundModel := builder.registeredModels[currConfig.Type]

		if foundModel {
			configFunction := builder.registeredModels[currConfig.Type]
			modelList[index] = configFunction(currConfig)

			parameterisedExplorer, ok := modelList[index].(parameters.Container)
			if ok {
				if errors := parameterisedExplorer.ParameterErrors(); errors != nil {
					wrappedErrors := errors2.Wrap(errors, "building model ["+currConfig.Name+"]")
					builder.errors.Add(wrappedErrors)
				}
			}
		} else {
			builder.errors.Add(
				errors.New("configuration specifies a model type [\"" +
					currConfig.Type + "\"], but no models are registered for that type"),
			)
		}
	}
	return modelList
}
