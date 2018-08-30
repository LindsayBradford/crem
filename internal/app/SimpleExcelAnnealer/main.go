// +build windows

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"
	"runtime"

	"github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/SimpleExcelAnnealer/components"
	"github.com/LindsayBradford/crm/logging/handlers"
	"github.com/pkg/errors"
)

var (
	defaultLogHandler   handlers.LogHandler
	oleFunctionsChannel = make(chan func())
	oleFunctionHandler  = func() {
		for oleFunction := range oleFunctionsChannel {
			oleFunction()
		}
	}
)

func doOleFunction(f func()) {
	done := make(chan bool, 1)
	oleFunctionsChannel <- func() {
		f()
		done <- true
	}
	<-done
}

func closeOleFunctionChannel() {
	close(oleFunctionsChannel)
	runtime.UnlockOSThread()
}

func main() {
	args := commandline.ParseArguments()
	RunFromConfigFile(args.ConfigFile)
}

func RunFromConfigFile(configFile string) {
	runtime.LockOSThread()
	scenarioConfig := retrieveConfig(configFile)
	scenarioRunner := buildScenarioOffConfig(scenarioConfig)
	runScenario(scenarioRunner)
}

func retrieveConfig(configFile string) *config.CRMConfig {
	configuration, retrieveError := config.Retrieve(configFile)

	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving simple excel annealer configuration")
		panic(wrappingError)
	}

	return configuration
}

func buildScenarioOffConfig(scenarioConfig *config.CRMConfig) annealing.CallableScenarioRunner {
	scenarioRunner, annealerLogHandler := components.BuildScenarioRunner(scenarioConfig, doOleFunction, closeOleFunctionChannel)
	defaultLogHandler = annealerLogHandler
	defaultLogHandler.Info("Configuring with [" + scenarioConfig.FilePath + "]")
	return scenarioRunner
}

func runScenario(scenarioRunner annealing.CallableScenarioRunner) {
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
	oleFunctionHandler()

	flushStreams()
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
