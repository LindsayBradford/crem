// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/pkg/excel"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

var (
	ScenarioLogger logging.Logger = loggers.DefaultNullLogger
)

func RunScenarioFromConfigFile(configFile string) {
	excel.EnableSpreadsheetSafeties()
	defer excel.DisableSpreadsheetSafeties()

	go runScenarioFromConfigFile(configFile)
	threading.GetMainThreadChannel().RunHandler()
}

func runScenarioFromConfigFile(configFile string) {
	configuration := retrieveScenarioConfiguration(configFile)
	establishScenarioLogger(configuration)
	scenario.RunScenarioFromConfig(configuration)
	threading.GetMainThreadChannel().Close()
}

func establishScenarioLogger(configuration *config.CREMConfig) {
	loggers, _ := new(config.LogHandlersBuilder).WithConfig(configuration.Loggers).Build()
	ScenarioLogger = loggers[0]
}

func retrieveScenarioConfiguration(configFile string) *config.CREMConfig {
	configuration, retrieveError := config.RetrieveCremFromFile(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		panic(wrappingError)
	}

	ScenarioLogger.Info("Configuring with [" + configuration.FilePath + "]")
	return configuration
}
