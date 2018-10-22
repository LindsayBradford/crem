// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"errors"

	"github.com/LindsayBradford/crem/annealing/solution"
	. "github.com/LindsayBradford/crem/errors"
)

type solutionExplorerBuilder struct {
	errors              *CompositeError
	config              []SolutionExplorerConfig
	registeredExplorers map[string]ExplorerConfigFunction
}

type ExplorerConfigFunction func(config SolutionExplorerConfig) solution.Explorer

func (builder *solutionExplorerBuilder) initialise() *solutionExplorerBuilder {
	if builder.errors == nil {
		builder.errors = new(CompositeError)
	}

	if builder.registeredExplorers == nil {
		builder.registerBaseExplorers()
	}

	return builder
}

func (builder *solutionExplorerBuilder) registerBaseExplorers() {
	builder.registeredExplorers = make(map[string]ExplorerConfigFunction, 2)

	builder.RegisteringExplorer(
		"NullSolutionExplorer",
		func(config SolutionExplorerConfig) solution.Explorer {
			return new(solution.NullExplorer).WithName(config.Name)
		},
	)

	builder.RegisteringExplorer(
		"DumbSolutionExplorer",
		func(config SolutionExplorerConfig) solution.Explorer {
			return new(solution.DumbExplorer).WithName(config.Name)
		},
	)
}

func (builder *solutionExplorerBuilder) WithConfig(cremConfig *CREMConfig) *solutionExplorerBuilder {
	builder.initialise()
	builder.config = cremConfig.SolutionExplorers
	return builder
}

func (builder *solutionExplorerBuilder) RegisteringExplorer(explorerType string, configFunction ExplorerConfigFunction) *solutionExplorerBuilder {
	builder.initialise()
	builder.registeredExplorers[explorerType] = configFunction
	return builder
}

func (builder *solutionExplorerBuilder) Build(explorerName string) (solution.Explorer, error) {
	var myExplorer solution.Explorer
	if len(builder.config) == 0 {
		builder.errors.Add(errors.New("configuration failed to specify any solution explorers"))
	} else {
		myExplorer = builder.findMyExplorer(explorerName, builder.buildExplorers())
	}

	if builder.errors.Size() > 0 {
		return nil, builder.errors
	}

	return myExplorer, nil
}

func (builder *solutionExplorerBuilder) findMyExplorer(myExplorerName string, explorers []solution.Explorer) solution.Explorer {
	for _, explorer := range explorers {
		if explorer != nil && explorer.Name() == myExplorerName {
			return explorer
		}
	}
	builder.errors.Add(
		errors.New("configuration specifies a non-existent solution explorer [\"" +
			myExplorerName + "\"] for its Annealer"),
	)
	return nil
}

func (builder *solutionExplorerBuilder) buildExplorers() []solution.Explorer {
	explorerList := make([]solution.Explorer, len(builder.config))
	for index, currConfig := range builder.config {
		_, foundExplorer := builder.registeredExplorers[currConfig.Type]

		if foundExplorer {
			configFunction := builder.registeredExplorers[currConfig.Type]
			explorerList[index] = configFunction(currConfig)
		} else {
			builder.errors.Add(
				errors.New("configuration specifies a solution explorer type [\"" +
					currConfig.Type + "\"], but no explorers are registered for that type"),
			)
		}
	}
	return explorerList
}
