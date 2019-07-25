// Copyright (c) 2019 Australian Rivers Institute.

package bootstrap

import (
	"os"

	"github.com/LindsayBradford/crem/cmd/cremexplorer/commandline"
	"github.com/LindsayBradford/crem/internal/pkg/config/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/config/userconfig/interpreter"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/excel"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

var (
	LogHandler    logging.Logger
	myScenario    scenario.Scenario
	myInterpreter interpreter.ConfigInterpreter
)

func init() {
	myInterpreter = *interpreter.NewInterpreter()
}

func RunExcelCompatibleScenarioFromConfigFile(configFile string) {
	excel.EnableSpreadsheetSafeties()
	defer excel.DisableSpreadsheetSafeties()

	go runMainThreadBoundScenarioFromConfigFile(configFile)
	threading.GetMainThreadChannel().RunHandler()
}

func runMainThreadBoundScenarioFromConfigFile(configFile string) {
	RunScenarioFromConfigFile(configFile)
	threading.GetMainThreadChannel().Close()
}

func RunScenarioFromConfigFile(configFile string) {
	deriveScenario(configFile)
	runScenario()
	flushStreams()
}

func runScenario() {
	if runError := myScenario.Run(); runError != nil {
		wrappingError := errors.Wrap(runError, "running scenario")
		LogHandler.Error(wrappingError)
		commandline.Exit(runError)
	}
}

func flushStreams() {
	os.Stdout.Sync()
	os.Stderr.Sync()
}

func deriveScenario(configFile string) {
	config := loadScenarioConfig(configFile)
	myScenario = myInterpreter.Interpret(config).Scenario()

	LogHandler = myScenario.LogHandler()
	LogHandler.Info("Configuring scenario with [" + config.MetaData.FilePath + "]")

	interpreterErrors := myInterpreter.Errors()

	if interpreterErrors != nil {
		wrappingError := errors.Wrap(interpreterErrors, "interpreting scenario file")
		commandline.Exit(wrappingError)
	}
}

func loadScenarioConfig(configFile string) *data.Config {
	configuration, retrieveError := data.RetrieveConfigFromFile(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		commandline.Exit(wrappingError)
	}

	return configuration
}
