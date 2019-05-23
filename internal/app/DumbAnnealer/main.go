// +build windows
// (c) 2018 Australian Rivers Institute.

package main

import (
	"os"

	"github.com/LindsayBradford/crem/internal/pkg/commandline"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

var (
	defaultLogHandler logging.Logger
)

func main() {
	args := commandline.ParseArguments()
	RunFromConfigFile(args.ScenarioFile)
}

func RunFromConfigFile(configFile string) {
	scenarioConfig := retrieveConfig(configFile)
	scenarioRunner := buildScenarioRunner(scenarioConfig)
	runScenario(scenarioRunner)
}

func retrieveConfig(configFile string) *config.CREMConfig {
	configuration, retrieveError := config.RetrieveCremFromFile(configFile)

	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving dumb annealer configuration")
		panic(wrappingError)
	}

	return configuration
}

func buildScenarioRunner(scenarioConfig *config.CREMConfig) scenario.CallableRunner {
	scenarioRunner, annealerLogHandler := buildRunnerAndLogger(scenarioConfig)

	defaultLogHandler = annealerLogHandler
	defaultLogHandler.Info("Configuring with [" + scenarioConfig.FilePath + "]")

	return scenarioRunner
}

func buildRunnerAndLogger(scenarioConfig *config.CREMConfig) (scenario.CallableRunner, logging.Logger) {
	newAnnealer, humanLogHandler, buildError :=
		new(config.AnnealerBuilder).
			WithConfig(scenarioConfig).
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
		WithTearDownFunction(threading.GetMainThreadChannel().Close).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber)

	if scenarioConfig.CpuProfilePath != "" {
		profilableRunner := new(scenario.ProfilableRunner).
			ThatProfiles(runner).
			ToFile(scenarioConfig.CpuProfilePath)

		runner = profilableRunner
	}

	return runner, humanLogHandler
}

func runScenario(scenarioRunner scenario.CallableRunner) {
	if runError := scenarioRunner.Run(); runError != nil {
		wrappingError := errors.Wrap(runError, "running dumb annealer scenario")
		defaultLogHandler.Error(wrappingError)
	}
	flushStreams()
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
