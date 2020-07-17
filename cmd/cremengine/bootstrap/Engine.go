// Copyright (c) 2019 Australian Rivers Institute.

package bootstrap

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/config/data"
	"github.com/LindsayBradford/crem/cmd/cremengine/config/interpreter"
	"github.com/LindsayBradford/crem/cmd/cremengine/engine"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"os"

	"github.com/LindsayBradford/crem/cmd/cremexplorer/commandline"
	"github.com/LindsayBradford/crem/pkg/excel"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

var (
	LogHandler    logging.Logger
	myEngine      engine.Engine
	myInterpreter interpreter.EngineConfigInterpreter
)

func init() {
	myInterpreter = *interpreter.NewEngineConfigInterpreter()
	LogHandler = loggers.DefaultTestingLogger // TODO: get final log handler wired in.
}

func RunMainThreadBoundEngineFromConfigFile(configFile string) {
	defer gracefullyHandlePanics()

	excel.EnableSpreadsheetSafeties()
	defer excel.DisableSpreadsheetSafeties()

	go runMainThreadBoundEngineFromConfigFile(configFile)
	threading.GetMainThreadChannel().RunHandler()
}

func gracefullyHandlePanics() {
	if r := recover(); r != nil {
		if recoveredError, isError := r.(error); isError {
			wrappingError := errors.Wrap(recoveredError, "running excel-compatible scenario")
			LogHandler.Error(wrappingError)
		}
		commandline.Exit(r)
	}
}

func runMainThreadBoundEngineFromConfigFile(configFile string) {
	defer func() {
		if r := recover(); r != nil {
			if recoveredError, isError := r.(error); isError {
				wrappingError := errors.Wrap(recoveredError, "running main-thread bound scenario")
				LogHandler.Error(wrappingError)
			}
			commandline.Exit(r)
		}
	}()

	RunEngineFromConfigFile(configFile)
	defer threading.GetMainThreadChannel().Close()
}

func RunEngineFromConfigFile(configFile string) {
	deriveEngineBehaviour(configFile)
	runEngine()
	flushStreams()
}

func runEngine() {
	if runError := myEngine.Run(); runError != nil {
		wrappingError := errors.Wrap(runError, "running engine")
		LogHandler.Error(wrappingError)
		commandline.Exit(runError)
	}
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}

func deriveEngineBehaviour(configFile string) {
	myConfig := loadConfig(configFile)
	myEngine = myInterpreter.Interpret(myConfig.Engine).Engine()

	myEngine.LogHandler().Info("Configuring with [" + myConfig.MetaData.FilePath + "]")

	interpreterErrors := myInterpreter.Errors()

	if interpreterErrors != nil {
		wrappingError := errors.Wrap(interpreterErrors, "interpreting engine configuration file")
		commandline.Exit(wrappingError)
	}
}

func loadConfig(configFile string) *data.EngineConfig {
	configuration, retrieveError := data.RetrieveConfigFromFile(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving engine configuration")
		commandline.Exit(wrappingError)
	}

	return configuration
}
