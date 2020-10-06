// Copyright (c) 2019 Australian Rivers Institute.

package bootstrap

import (
	"fmt"
	"os"

	"github.com/LindsayBradford/crem/cmd/cremexplorer/commandline"
	data2 "github.com/LindsayBradford/crem/cmd/cremexplorer/config/data"
	interpreter2 "github.com/LindsayBradford/crem/cmd/cremexplorer/config/interpreter"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/excel"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

var (
	LogHandler    logging.Logger
	myScenario    scenario.Scenario
	myInterpreter interpreter2.ConfigInterpreter
)

func init() {
	myInterpreter = *interpreter2.NewInterpreter()
}

func RunExcelCompatibleScenarioFromConfigFile(configFile string) {
	defer gracefullyHandlePanics()

	excel.EnableSpreadsheetSafeties()
	defer excel.DisableSpreadsheetSafeties()

	go runMainThreadBoundScenarioFromConfigFile(configFile)
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

func runMainThreadBoundScenarioFromConfigFile(configFile string) {
	defer func() {
		if r := recover(); r != nil {
			if recoveredError, isError := r.(error); isError {
				wrappingError := errors.Wrap(recoveredError, "running main-thread bound scenario")
				LogHandler.Error(wrappingError)
			}
			commandline.Exit(r)
		}
	}()

	RunScenarioFromConfigFile(configFile)
	defer threading.GetMainThreadChannel().Close()
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
	myConfig := loadScenarioConfig(configFile)
	myScenario = myInterpreter.Interpret(myConfig).Scenario()

	LogHandler = myScenario.LogHandler()
	metaData := myConfig.MetaData
	metaDataSummary := fmt.Sprintf("Running [%s] Version [%s] with scenario [%s]",
		metaData.ExecutableName, metaData.ExecutableVersion, metaData.FilePath)
	LogHandler.Info(metaDataSummary)

	interpreterErrors := myInterpreter.Errors()

	if interpreterErrors != nil {
		wrappingError := errors.Wrap(interpreterErrors, "interpreting scenario file")
		commandline.Exit(wrappingError)
	}
}

func loadScenarioConfig(configFile string) *data2.Config {
	configuration, retrieveError := data2.RetrieveConfigFromFile(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		commandline.Exit(wrappingError)
	}

	return configuration
}
