package interpreter

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding"
	"github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/formatters"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

type ScenarioConfigInterpreter struct {
	errors *compositeErrors.CompositeError

	scenario scenario.Scenario
}

func NewScenarioConfigInterpreter() *ScenarioConfigInterpreter {
	interpreter := new(ScenarioConfigInterpreter).initialise()
	return interpreter
}

func (i *ScenarioConfigInterpreter) initialise() *ScenarioConfigInterpreter {
	i.errors = compositeErrors.New("Scenario Configuration")
	i.scenario = scenario.NullScenario
	return i
}

func (i *ScenarioConfigInterpreter) Interpret(scenarioConfig *data.ScenarioConfig) *ScenarioConfigInterpreter {
	scenario := scenario.NewBaseScenario()

	runner := buildRunner(scenarioConfig)
	scenario.SetRunner(runner)

	i.scenario = scenario
	return i
}

func buildRunner(scenarioConfig *data.ScenarioConfig) scenario.CallableRunner {
	var runner scenario.CallableRunner

	saver := buildSaver(scenarioConfig)
	logHandler := buildLogHandler(scenarioConfig)

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

func buildLogHandler(scenarioConfig *data.ScenarioConfig) logging.Logger {
	// TODO: build off config.
	builder := new(loggers.Builder)

	defaultLogger, _ := builder.ForNativeLibraryLogHandler().
		WithName("DefaultLogHandler").
		WithFormatter(new(formatters.RawMessageFormatter)).
		WithLogLevelDestination(observer.AnnealerLogLevel, logging.STDOUT).
		WithLogLevelDestination(model.LogLevel, logging.DISCARD).
		Build()

	return defaultLogger
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
