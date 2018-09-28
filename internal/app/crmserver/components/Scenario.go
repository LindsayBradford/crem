// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/crmserver/components/scenario"
	"github.com/LindsayBradford/crm/logging/handlers"
	"github.com/pkg/errors"
)

var (
	ScenarioLogger handlers.LogHandler = handlers.DefaultNullLogHandler
)

func RunScenarioFromConfigFile(configFile string) {
	configuration := retrieveScenarioConfiguration(configFile)
	establishScenarioLogger(configuration)
	scenario.RunScenarioFromConfig(configuration)
}

func establishScenarioLogger(configuration *config.CRMConfig) {
	loggers, _ := new(config.LogHandlersBuilder).WithConfig(configuration.Loggers).Build()
	ScenarioLogger = loggers[0]
}

func retrieveScenarioConfiguration(configFile string) *config.CRMConfig {
	configuration, retrieveError := config.RetrieveCrmFromFile(configFile)
	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		panic(wrappingError)
	}

	ScenarioLogger.Info("Configuring with [" + configuration.FilePath + "]")
	return configuration
}
