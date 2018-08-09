// Copyright (c) 2018 Australian Rivers Institute.

package config

import (
	"errors"

	"github.com/LindsayBradford/crm/annealing/solution"
	. "github.com/LindsayBradford/crm/errors"
)

type SolutionExplorerBuilder struct {
	errors              *CompositeError
	config              []SolutionExplorerConfig
	registeredExplorers map[string]ExplorerConfigFunction
}

type ExplorerConfigFunction func(config SolutionExplorerConfig) solution.SolutionExplorer

func (builder *SolutionExplorerBuilder) initialise() *SolutionExplorerBuilder {
	if builder.errors == nil {
		builder.errors = new(CompositeError)
	}

	if builder.registeredExplorers == nil {
		builder.registerBaseExplorers()
	}

	return builder
}

func (builder *SolutionExplorerBuilder) registerBaseExplorers() {
	builder.registeredExplorers = make(map[string]ExplorerConfigFunction, 2)

	builder.RegisteringExplorer(
		"NullSolutionExplorer",
		func(config SolutionExplorerConfig) solution.SolutionExplorer {
			return new(solution.NullSolutionExplorer).WithName(config.Name)
		},
	)

	builder.RegisteringExplorer(
		"DumbSolutionExplorer",
		func(config SolutionExplorerConfig) solution.SolutionExplorer {
			return new(solution.DumbSolutionExplorer).WithName(config.Name)
		},
	)
}

func (builder *SolutionExplorerBuilder) WithConfig(crmConfig *CRMConfig) *SolutionExplorerBuilder {
	builder.initialise()
	builder.config = crmConfig.SolutionExplorers
	return builder
}

func (builder *SolutionExplorerBuilder) RegisteringExplorer(explorerName string, configFunction ExplorerConfigFunction) *SolutionExplorerBuilder {
	builder.initialise()
	builder.registeredExplorers[explorerName] = configFunction
	return builder
}

func (builder *SolutionExplorerBuilder) Build(explorerName string) (solution.SolutionExplorer, error) {
	var myExplorer solution.SolutionExplorer
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

func (builder *SolutionExplorerBuilder) findMyExplorer(myExplorerName string, explorers []solution.SolutionExplorer) solution.SolutionExplorer {
	for _, explorer := range explorers {
		if explorer.Name() == myExplorerName {
			return explorer
		}
	}
	builder.errors.Add(
		errors.New("configuration specifies a non-existent solution explorer [\"" +
			myExplorerName + "\"] for its Annealer"),
	)
	return nil
}

func (builder *SolutionExplorerBuilder) buildExplorers() []solution.SolutionExplorer {
	explorerList := make([]solution.SolutionExplorer, len(builder.config))
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
