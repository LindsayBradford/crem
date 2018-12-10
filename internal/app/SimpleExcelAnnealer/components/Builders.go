// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"os"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/threading"
)

func BuildScenarioRunner(scenarioConfig *config.CREMConfig, mainThreadChannel *threading.MainThreadChannel) (scenario.CallableRunner, logging.Logger) {
	newAnnealer, humanLogHandler, buildError :=
		new(config.AnnealerBuilder).
			WithConfig(scenarioConfig).
			RegisteringModel(buildSimpleExcelModelRegistration(mainThreadChannel.Call)).
			RegisteringExplorer(buildSimpleExcelExplorerRegistration()).
			Build()

	if buildError != nil {
		humanLogHandler.Error(buildError)
		humanLogHandler.Error("Exiting program due to failed Annealer build")
		os.Exit(1)
	}

	var runner scenario.CallableRunner

	runner = new(scenario.Runner).
		ForAnnealer(newAnnealer).
		WithName(scenarioConfig.ScenarioName).
		WithRunNumber(scenarioConfig.RunNumber).
		WithTearDownFunction(mainThreadChannel.Close).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber)

	runner = new(scenario.SpreadsheetSafeScenarioRunner).ThatLocks(runner)

	if scenarioConfig.CpuProfilePath != "" {
		profilableRunner := new(scenario.ProfilableRunner).
			ThatProfiles(runner).
			ToFile(scenarioConfig.CpuProfilePath)

		runner = profilableRunner
	}

	return runner, humanLogHandler
}

func buildSimpleExcelModelRegistration(wrapper threading.MainThreadFunctionWrapper) config.ModelRegistration {
	return config.ModelRegistration{
		ModelType: "SimpleExcelModel",
		ConfigFunction: func(config config.ModelConfig) model.Model {
			return NewSimpleExcelModel().
				WithParameters(config.Parameters).
				WithName(config.Name).
				WithOleFunctionWrapper(wrapper)
		},
	}
}

func buildSimpleExcelExplorerRegistration() config.ExplorerRegistration {
	return config.ExplorerRegistration{
		ExplorerType: "SimpleExcelExplorer",
		ConfigFunction: func(config config.SolutionExplorerConfig) explorer.Explorer {
			return NewSimpleExcelExplorer().
				WithParameters(config.Parameters).
				WithName(config.Name)
		},
	}
}
