// +build windows

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package main

import (
	"os"

	"github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/commandline"
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/SimpleExcelAnnealer/components"
	"github.com/LindsayBradford/crm/logging/handlers"
	"github.com/LindsayBradford/crm/profiling"
	"github.com/pkg/errors"
)

var (
	args               = commandline.ParseArguments()
	annealingFunctions = new(profiling.OptionalProfilingFunctionPair)
	defaultLogHandler  handlers.LogHandler
)

func main() {
	buildAnnealingFunctions()

	var runErr error
	if profilingRequested() {
		runErr = annealingFunctions.ProfiledFunction()
	} else {
		runErr = annealingFunctions.UnProfiledFunction()
	}

	if runErr != nil {
		os.Exit(1)
	}

	defer flushStreams()
}

func buildAnnealingFunctions() {
	annealer := buildAnnealerOffConfig()

	handleRecovery := func(wrapperMessage string, recoveryMsg interface{}) error {
		recoveryError, ok := recoveryMsg.(error)
		if ok {
			wrappedError := errors.Wrap(recoveryError, wrapperMessage)
			defaultLogHandler.Error(wrappedError)
			return wrappedError
		}
		return nil
	}

	annealingFunctions.UnProfiledFunction = func() (functionError error) {
		defer func() {
			if r := recover(); r != nil {
				functionError = handleRecovery("un-profiled annealer", r)
			}
		}()

		defaultLogHandler.Debug("About to call annealer.Anneal()")
		annealer.Anneal()
		defaultLogHandler.Debug("Call to annealer.Anneal() finished.")
		return
	}

	annealingFunctions.ProfiledFunction = func() (functionError error) {
		defer func() {
			if r := recover(); r != nil {
				functionError = handleRecovery("profiled annealer", r)
			}
		}()

		defer defaultLogHandler.Debug("Cpu profiling to file [" + args.CpuProfile + "] now generated")
		return profiling.CpuProfileOfFunctionToFile(annealingFunctions.UnProfiledFunction, args.CpuProfile)
	}
}

func buildAnnealerOffConfig() shared.Annealer {
	config := retrieveConfig()
	annealer, annealerLogHandler := components.BuildAnnealer(config)
	defaultLogHandler = annealerLogHandler
	return annealer
}

func retrieveConfig() *config.CRMConfig {
	configuration, retrieveError := config.Retrieve(args.ConfigFile)
	if retrieveError != nil {
		commandline.Exit(retrieveError)
	}
	return configuration
}

func profilingRequested() bool {
	return args.CpuProfile != ""
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}
