package interpreter

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding"
	"github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
)

type ScenarioConfigInterpreter struct {
	errors *compositeErrors.CompositeError

	observerInterpreter *ObserverConfigInterpreter

	scenario scenario.Scenario
}

func NewScenarioConfigInterpreter() *ScenarioConfigInterpreter {
	interpreter := new(ScenarioConfigInterpreter).initialise()
	return interpreter
}

func (i *ScenarioConfigInterpreter) initialise() *ScenarioConfigInterpreter {
	i.errors = compositeErrors.New("Scenario Configuration")
	i.scenario = scenario.NullScenario
	i.observerInterpreter = NewObserverConfigInterpreter()
	return i
}

func (i *ScenarioConfigInterpreter) Interpret(scenarioConfig *data.ScenarioConfig) *ScenarioConfigInterpreter {
	i.observerInterpreter.Interpret(&scenarioConfig.Observer)

	runner := i.buildRunner(scenarioConfig)
	i.scenario = scenario.NewBaseScenario().
		WithRunner(runner).
		WithObserver(i.observerInterpreter.Observer())

	return i
}

func (i *ScenarioConfigInterpreter) buildRunner(scenarioConfig *data.ScenarioConfig) scenario.CallableRunner {
	var runner scenario.CallableRunner

	saver := buildSaver(scenarioConfig)
	logHandler := i.observerInterpreter.LogHandler()

	runner = scenario.NewRunner().
		WithName(scenarioConfig.Name).
		WithRunNumber(scenarioConfig.RunNumber).
		WithMaximumConcurrentRuns(scenarioConfig.MaximumConcurrentRunNumber).
		WithLogHandler(logHandler).
		WithSaver(saver)

	if scenarioConfig.CpuProfilePath != "" {
		profilingRunner := new(scenario.ProfilingRunner).
			ThatProfiles(runner).
			ToFile(scenarioConfig.CpuProfilePath)

		runner = profilingRunner
	}

	return runner
}

func buildSaver(scenarioConfig *data.ScenarioConfig) *scenario.Saver {
	saver := scenario.NewSaver().
		WithOutputType(configOutputTypeToSolutionOutputType(scenarioConfig.OutputType)).
		WithOutputPath(scenarioConfig.OutputPath)
	return saver
}

func configOutputTypeToSolutionOutputType(outputType data.ScenarioOutputType) encoding.OutputType {
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
