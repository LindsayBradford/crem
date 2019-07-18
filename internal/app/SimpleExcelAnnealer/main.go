// +build windows

// (c) 2018 Australian Rivers Institute.
package main

import (
	"os"

	"github.com/LindsayBradford/crem/cmd/cremexplorer/commandline"
	"github.com/LindsayBradford/crem/internal/app/SimpleExcelAnnealer/components"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

var (
	defaultLogHandler logging.Logger
	mainThreadChannel threading.MainThreadChannel
)

func init() {
	mainThreadChannel = threading.GetMainThreadChannel()
}

func main() {
	args := commandline.ParseArguments()
	RunFromConfigFile(args.ScenarioFile)
}

func RunFromConfigFile(configFile string) {
	scenarioConfig := retrieveConfig(configFile)
	scenarioRunner := buildScenarioOffConfig(scenarioConfig)
	runScenario(scenarioRunner)
}

func retrieveConfig(configFile string) *config.CREMConfig {
	configuration, retrieveError := config.RetrieveCremFromFile(configFile)

	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving simple excel annealer configuration")
		panic(wrappingError)
	}

	return configuration
}

func buildScenarioOffConfig(scenarioConfig *config.CREMConfig) scenario.CallableRunner {
	scenarioRunner, annealerLogHandler := components.BuildScenarioRunner(scenarioConfig, &mainThreadChannel)
	defaultLogHandler = annealerLogHandler
	defaultLogHandler.Info("Configuring with [" + scenarioConfig.FilePath + "]")
	return scenarioRunner
}

func runScenario(scenarioRunner scenario.CallableRunner) {
	defer func() {
		if r := recover(); r != nil {
			recoveryError, ok := r.(error)
			if ok {
				wrappedError := errors.Wrap(recoveryError, "running simple excel annealer scenario")
				defaultLogHandler.Error(wrappedError)
				panic(wrappedError)
			}
		}
	}()

	defaultLogHandler.Debug("starting scenario runner")
	go scenarioRunner.Run()

	defaultLogHandler.Debug("starting ole function polling")
	mainThreadChannel.RunHandler()

	flushStreams()
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
