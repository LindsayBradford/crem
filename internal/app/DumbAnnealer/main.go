// +build windows
// (c) 2018 Australian Rivers Institute.

package main

import (
	"os"

	"github.com/LindsayBradford/crem/cmd/cremexplorer/commandline"
	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/config/interpreter"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
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

func main() {
	args := commandline.ParseArguments()

	go RunFromConfigFile(args.ScenarioFile)
	threading.GetMainThreadChannel().RunHandler()
}

func RunFromConfigFile(configFile string) {
	deriveScenario(configFile)
	runScenario()
	flushStreams()
	threading.GetMainThreadChannel().Close()
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
