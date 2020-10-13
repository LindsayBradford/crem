// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	appData "github.com/LindsayBradford/crem/cmd/cremexplorer/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/pkg/errors"
)

type ScenarioConfigInterpreter struct {
	errors *compositeErrors.CompositeError

	reportingInterpreter *ReportingConfigInterpreter

	scenario scenario.Scenario
	runner   scenario.CallableRunner
}

func NewScenarioConfigInterpreter() *ScenarioConfigInterpreter {
	interpreter := new(ScenarioConfigInterpreter).initialise()
	return interpreter
}

func (i *ScenarioConfigInterpreter) initialise() *ScenarioConfigInterpreter {
	i.errors = compositeErrors.New("Scenario Configuration")
	i.scenario = scenario.NullScenario
	i.runner = scenario.NullRunner
	i.reportingInterpreter = NewObserverConfigInterpreter()

	return i
}

func (i *ScenarioConfigInterpreter) Interpret(scenarioConfig *appData.ScenarioConfig) *ScenarioConfigInterpreter {
	i.interpretReporting(&scenarioConfig.Reporting)
	i.interpretRunner(scenarioConfig)

	i.scenario = scenario.NewBaseScenario().
		WithRunner(i.runner).
		WithObserver(i.reportingInterpreter.Observer())

	return i
}

func (i *ScenarioConfigInterpreter) interpretReporting(config *appData.ReportingConfig) {
	i.reportingInterpreter.Interpret(config)
	if i.reportingInterpreter.Errors() != nil {
		i.errors.Add(i.reportingInterpreter.Errors())
	}
}

func (i *ScenarioConfigInterpreter) interpretRunner(config *appData.ScenarioConfig) {
	var runner scenario.CallableRunner

	if config.Name == "" {
		i.errors.Add(errors.New("Missing mandatory scenario name field"))
	}

	logHandler := i.reportingInterpreter.LogHandler()
	saver := buildSaver(config).WithLogHandler(logHandler)

	runner = scenario.NewRunner().
		WithName(config.Name).
		WithRunNumber(config.RunNumber).
		WithMaximumConcurrentRuns(config.MaximumConcurrentRunNumber).
		WithLogHandler(logHandler).
		WithSaver(saver)

	if config.CpuProfilePath != "" {
		profilingRunner := new(scenario.ProfilingRunner).
			ThatProfiles(runner).
			ToFile(config.CpuProfilePath)

		runner = profilingRunner
	}

	i.runner = runner
}

func buildSaver(scenarioConfig *appData.ScenarioConfig) *scenario.Saver {
	saver := scenario.NewSaver().
		WithOutputType(configOutputTypeToSolutionOutputType(scenarioConfig.OutputType)).
		WithOutputPath(scenarioConfig.OutputPath)
	return saver
}

func configOutputTypeToSolutionOutputType(outputType appData.ScenarioOutputType) encoding.OutputType {
	return encoding.OutputType(outputType.String())
}

func (i *ScenarioConfigInterpreter) Scenario() scenario.Scenario {
	return i.scenario
}

func (i *ScenarioConfigInterpreter) Errors() error {
	if i.errors.Size() > 0 {
		return i.errors
	}
	return nil
}
