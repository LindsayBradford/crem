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

func init() {
	// Arrange that main.main runs on main thread.
	runtime.LockOSThread()
}

var (
	args                = commandline.ParseArguments()
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
	scenarioConfig := retrieveConfig()

	defer func() {
		if r := recover(); r != nil {
			recoveryError, ok := r.(error)
			if ok {
				wrappedError := errors.Wrap(recoveryError, "main")
				defaultLogHandler.Error(wrappedError)
				commandline.Exit(1)
			}
		}
	}()

	scenarioRunner := buildScenarioOffConfig(scenarioConfig)

	defaultLogHandler.Debug("starting scenario runner")
	go scenarioRunner.Run()

	defaultLogHandler.Debug("starting ole function polling")
	oleFunctionHandler()

	defer flushStreams()
}

func buildScenarioOffConfig(scenarioConfig *config.CRMConfig) annealing.CallableScenarioRunner {
	scenarioRunner, annealerLogHandler := components.BuildScenarioRunner(scenarioConfig, doOleFunction, closeOleFunctionChannel)
	defaultLogHandler = annealerLogHandler
	return scenarioRunner
}

func retrieveConfig() *config.CRMConfig {
	configuration, retrieveError := config.Retrieve(args.ConfigFile)
	if retrieveError != nil {
		commandline.Exit(retrieveError)
	}
	return configuration
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
