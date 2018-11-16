// +build windows

// (c) 2018 Australian Rivers Institute.
package main

import (
	"os"
	"runtime"

	"github.com/LindsayBradford/crem/internal/app/SimpleExcelAnnealer/components"
	"github.com/LindsayBradford/crem/internal/pkg/commandline"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/pkg/errors"
)

var (
	defaultLogHandler logging.Logger
	mainThreadChannel = make(chan func())
)

func init() {
	runtime.LockOSThread() // ensure main goroutine is locked to OS thread for proper callOnMainThread behaviour.
}

func mainThreadFunctionHandler() {
	for function := range mainThreadChannel {
		function()
	}
}

func callOnMainThread(function func()) {
	done := make(chan bool, 1)
	mainThreadChannel <- func() {
		function()
		done <- true
	}
	<-done
}

func closeMainThreadChannel() {
	close(mainThreadChannel)
	runtime.UnlockOSThread()
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
	scenarioRunner, annealerLogHandler := components.BuildScenarioRunner(scenarioConfig, callOnMainThread, closeMainThreadChannel)
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
	mainThreadFunctionHandler()

	flushStreams()
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
