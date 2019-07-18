// Copyright (c) 2019 Australian Rivers Institute.

package bootstrap

import (
	"github.com/LindsayBradford/crem/cmd/cremexplorer/commandline"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/pkg/excel"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

var (
	LogHandler logging.Logger
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
	establishScenarioLogger(configuration)

	if runError := RunScenarioFromConfig(configuration); runError != nil {
		commandline.Exit(runError)
	}
}

func establishScenarioLogger(configuration *config.CREMConfig) {
	loggers, _ := new(config.LogHandlersBuilder).WithConfig(configuration.Loggers).Build()
	LogHandler = loggers[0]
	LogHandler.Info("Configuring with [" + configuration.FilePath + "]")
}

func retrieveScenarioConfiguration(configFile string) *config.CREMConfig {
	configuration, retrieveError := config.RetrieveCremFromFile(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		panic(wrappingError)
	}

	return configuration
}
