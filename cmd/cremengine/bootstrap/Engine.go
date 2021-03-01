// Copyright (c) 2019 Australian Rivers Institute.

package bootstrap

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/config/data"
	"github.com/LindsayBradford/crem/cmd/cremengine/config/interpreter"
	"github.com/LindsayBradford/crem/cmd/cremengine/engine"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"os"

	"github.com/LindsayBradford/crem/cmd/cremengine/commandline"
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
	LogHandler = loggers.DefaultTestingLogger
}

func RunMainThreadBoundEngineFromArguments(args *commandline.Arguments) {
	defer gracefullyHandlePanics()

	excel.EnableSpreadsheetSafeties()
	defer excel.DisableSpreadsheetSafeties()

	go runMainThreadBoundEngineFromArguments(args)
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

func runMainThreadBoundEngineFromArguments(args *commandline.Arguments) {
	defer func() {
		if r := recover(); r != nil {
			if recoveredError, isError := r.(error); isError {
				wrappingError := errors.Wrap(recoveredError, "running main-thread bound scenario")
				LogHandler.Error(wrappingError)
			}
			commandline.Exit(r)
		}
	}()

	RunEngineFromArguments(args)
	defer threading.GetMainThreadChannel().Close()
}

func RunEngineFromArguments(args *commandline.Arguments) {
	deriveEngineBehaviour(args)
	deriveInitialEngineState(args)
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

func deriveEngineBehaviour(args *commandline.Arguments) {
	myConfig := loadConfig(args.EngineConfigFile)
	myEngine = myInterpreter.Interpret(myConfig.Engine).Engine()

	myEngine.LogHandler().Info("Configuring with [" + myConfig.MetaData.FilePath + "]")

	interpreterErrors := myInterpreter.Errors()

	if interpreterErrors != nil {
		wrappingError := errors.Wrap(interpreterErrors, "interpreting engine configuration file")
		commandline.Exit(wrappingError)
	}
}

func deriveInitialEngineState(args *commandline.Arguments) {
	if args.ScenarioFile != "" {
		myEngine.LogHandler().Info("Initialising engine with scenario [" + args.ScenarioFile + "]")

		myEngine.SetScenario(args.ScenarioFile)
	}
	if args.SolutionFile != "" {
		myEngine.LogHandler().Info("Initialising engine management actions with solution [" + args.SolutionFile + "]")

		myEngine.SetSolution(args.SolutionFile)
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
