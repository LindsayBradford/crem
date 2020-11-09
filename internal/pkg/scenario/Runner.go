// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"fmt"
	"sync"
	. "time"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging"
)

type CallableRunner interface {
	SetAnnealer(annealer annealing.Annealer)
	LogHandler() logging.Logger

	Run() error
}

type Runner struct {
	annealer   annealing.Annealer
	logHandler logging.Logger
	saver      CallableSaver

	name              string
	operationType     string
	runNumber         uint64
	maxConcurrentRuns uint64
	tearDown          func()

	startTime  Time
	finishTime Time
}

func NewRunner() *Runner {
	return new(Runner).initialise()
}

func (runner *Runner) initialise() *Runner {
	runner.runNumber = 1
	runner.maxConcurrentRuns = 1 // Sequential by default
	runner.name = "Default Scenario"
	runner.tearDown = defaultTearDown
	return runner
}

func (runner *Runner) WithLogHandler(logHandler logging.Logger) *Runner {
	runner.logHandler = logHandler
	return runner
}

func (runner *Runner) WithSaver(saver CallableSaver) *Runner {
	saver.SetLogHandler(runner.logHandler)
	runner.saver = saver
	return runner
}

func (runner *Runner) WithName(name string) *Runner {
	if name != "" {
		runner.name = name
	}
	return runner
}

func (runner *Runner) WithRunNumber(runNumber uint64) *Runner {
	if runNumber > 0 {
		runner.runNumber = runNumber
	}
	return runner
}

func (runner *Runner) WithTearDownFunction(tearDown func()) *Runner {
	if tearDown != nil {
		runner.tearDown = tearDown
	}
	return runner
}

func defaultTearDown() {
	// deliberately does nothing
}

func (runner *Runner) WithMaximumConcurrentRuns(maxConcurrentRuns uint64) *Runner {
	if maxConcurrentRuns > 0 {
		runner.maxConcurrentRuns = maxConcurrentRuns
	}
	return runner
}

func (runner *Runner) SetAnnealer(annealer annealing.Annealer) {
	runner.annealer = annealer
	annealer.SetLogHandler(runner.logHandler)
	runner.saver.SetDecompressionModel(annealer.Model())
	annealer.AddObserver(runner.saver)
}

func (runner *Runner) Run() error {
	runner.logScenarioStartMessage()
	runner.startTime = Now()

	runError := runner.runScenario()

	runner.finishTime = Now()
	runner.logHandler.Info("Finished running scenario [" + runner.name + "]")

	runner.logHandler.Info(runner.generateElapsedTimeString())

	runner.tearDown()

	return runError
}

func (runner *Runner) LogHandler() logging.Logger {
	return runner.logHandler
}

func (runner *Runner) logScenarioStartMessage() {
	var runTypeText string
	if runner.maxConcurrentRuns > 1 {
		runTypeText = fmt.Sprintf("executing a maximum of %d runs concurrently.", runner.maxConcurrentRuns)
	} else {
		runTypeText = "executing runs sequentially"
	}

	message := fmt.Sprintf("Scenario [%s]: configured for %d run(s), %s", runner.name, runner.runNumber, runTypeText)
	runner.logHandler.Info(message)
}

func (runner *Runner) generateElapsedTimeString() string {
	return fmt.Sprintf("Total elapsed time of scenario [%s] = [%v]", runner.name, runner.ElapsedTime())
}

func (runner *Runner) ElapsedTime() Duration {
	return runner.finishTime.Sub(runner.startTime)
}

func (runner *Runner) runScenario() error {
	var runWaitGroup sync.WaitGroup

	concurrentRunGuard := make(chan struct{}, runner.maxConcurrentRuns)

	doRun := func(runNumber uint64) {
		runner.run(runNumber)
		<-concurrentRunGuard
		runWaitGroup.Done()
	}

	runWaitGroup.Add(int(runner.runNumber))

	for runNumber := uint64(1); runNumber <= runner.runNumber; runNumber++ {
		concurrentRunGuard <- struct{}{}
		go doRun(runNumber)
	}

	runWaitGroup.Wait()

	return nil
}

func (runner *Runner) run(runNumber uint64) {
	annealerCopy := runner.annealer.DeepClone()

	runner.assignNewRunId(runNumber, annealerCopy)
	runner.wireObservers(annealerCopy)

	annealerCopy.Anneal()
	runner.logRunFinishedMessage(runNumber)
}

func (runner *Runner) assignNewRunId(runNumber uint64, annealerCopy annealing.Annealer) {
	runId := runner.generateCloneId(runNumber)
	annealerCopy.SetId(runId)
	annealerCopy.SolutionExplorer().SetId(runId)
	annealerCopy.SolutionExplorer().Model().SetId(runId)
	runner.logRunStartMessage(runNumber)
}

func (runner *Runner) wireObservers(annealer annealing.Annealer) {
	if observingAnnealer, annealerIsObserver := annealer.(observer.Observer); annealerIsObserver {
		explorer := annealer.SolutionExplorer()
		if eventNotifyingExplorer, explorerIssEventNotifier := explorer.(observer.EventNotifier); explorerIssEventNotifier {
			eventNotifyingExplorer.AddObserver(observingAnnealer)
		}
		model := annealer.Model()
		if eventNotifyingModel, modelIssEventNotifier := model.(observer.EventNotifier); modelIssEventNotifier {
			eventNotifyingModel.AddObserver(observingAnnealer)
		}
	}
}

func (runner *Runner) generateCloneId(runNumber uint64) string {
	if runner.runNumber > 1 {
		return fmt.Sprintf("%s (%d/%d)", runner.name, runNumber, runner.runNumber)
	} else {
		return runner.name
	}
}

func (runner *Runner) logRunStartMessage(runNumber uint64) {
	if runner.runNumber > 1 {
		runner.logHandler.Info(runner.generateCloneId(runNumber) + ": run started")
	}
}

func (runner *Runner) logRunFinishedMessage(runNumber uint64) {
	if runner.runNumber > 1 {
		runner.logHandler.Info(runner.generateCloneId(runNumber) + ": run finished")
	}
}

var NullRunner CallableRunner = new(nullRunner)

type nullRunner struct{}

func (s *nullRunner) SetAnnealer(annealer annealing.Annealer) {}
func (s *nullRunner) LogHandler() logging.Logger              { return nil }
func (s *nullRunner) Run() error                              { return nil }
