// Copyright (c) 2019 Australian Rivers Institute.

package bootstrap

import (
	"github.com/LindsayBradford/crem/cmd/cremexplorer/commandline"
	"github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/interpreter"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/excel"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

var (
	LogHandler logging.Logger
	Scenario   scenario.Scenario
)

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
	configuration := retrieveScenarioConfiguration(configFile)
	establishScenario(configuration)

	if runError := Scenario.Run(); runError != nil {
		commandline.Exit(runError)
	}
}

func establishScenario(config *data.Config) {
	Scenario = interpreter.NewInterpreter().Interpret(config).Scenario()

	LogHandler = Scenario.LogHandler()
	LogHandler.Info("Configuring with [" + config.MetaData.FilePath + "]")

}

func retrieveScenarioConfiguration(configFile string) *data.Config {
	configuration, retrieveError := data.RetrieveConfigFromFile(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		panic(wrappingError)
	}

	return configuration
}
