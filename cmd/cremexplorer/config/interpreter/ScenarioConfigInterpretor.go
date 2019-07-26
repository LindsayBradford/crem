// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	data2 "github.com/LindsayBradford/crem/cmd/cremexplorer/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
)

type ScenarioConfigInterpreter struct {
	errors *compositeErrors.CompositeError

	observerInterpreter *ReportingConfigInterpreter

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
	i.observerInterpreter = NewObserverConfigInterpreter()

	return i
}

func (i *ScenarioConfigInterpreter) Interpret(scenarioConfig *data2.ScenarioConfig) *ScenarioConfigInterpreter {
	i.interpretObserver(&scenarioConfig.Reporting)
	i.interpretRunner(scenarioConfig)

	i.scenario = scenario.NewBaseScenario().
		WithRunner(i.runner).
		WithObserver(i.observerInterpreter.Observer())

	return i
}

func (i *ScenarioConfigInterpreter) interpretObserver(config *data2.ReportingConfig) {
	i.observerInterpreter.Interpret(config)
	if i.observerInterpreter.Errors() != nil {
		i.errors.Add(i.observerInterpreter.Errors())
	}
}

func (i *ScenarioConfigInterpreter) interpretRunner(config *data2.ScenarioConfig) {
	var runner scenario.CallableRunner

	saver := buildSaver(config)
	logHandler := i.observerInterpreter.LogHandler()

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

func buildSaver(scenarioConfig *data2.ScenarioConfig) *scenario.Saver {
	saver := scenario.NewSaver().
		WithOutputType(configOutputTypeToSolutionOutputType(scenarioConfig.OutputType)).
		WithOutputPath(scenarioConfig.OutputPath)
	return saver
}

func configOutputTypeToSolutionOutputType(outputType data2.ScenarioOutputType) encoding.OutputType {
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
