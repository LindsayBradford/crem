// Copyright (c) 2018 Australian Rivers Institute.

package annealing

import (
	"fmt"
	"runtime"
	"sync"
	. "time"

	"github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/logging/handlers"
	"github.com/LindsayBradford/crm/profiling"
	"github.com/go-ole/go-ole"
)

type CallableScenarioRunner interface {
	Run() error
	LogHandler() handlers.LogHandler
}

type ScenarioRunner struct {
	annealer   shared.Annealer
	logHandler handlers.LogHandler

	name          string
	operationType string
	runNumber     uint64
	concurrently  bool
	tearDown      func()

	startTime  Time
	finishTime Time
}

func (runner *ScenarioRunner) ForAnnealer(annealer shared.Annealer) *ScenarioRunner {
	runner.runNumber = 1
	runner.concurrently = false
	runner.name = "Default Scenario"

	runner.logHandler = annealer.LogHandler()
	runner.annealer = annealer

	return runner
}

func (runner *ScenarioRunner) WithName(name string) *ScenarioRunner {
	if name == "" {
		return runner
	}
	runner.name = name
	return runner
}

func (runner *ScenarioRunner) WithRunNumber(runNumber uint64) *ScenarioRunner {
	if runNumber == 0 {
		return runner
	}
	runner.runNumber = runNumber
	return runner
}

func (runner *ScenarioRunner) WithTearDownFunction(tearDown func()) *ScenarioRunner {
	if tearDown != nil {
		runner.tearDown = tearDown
	} else {
		runner.tearDown = defaultTeadDown
	}
	return runner
}

func defaultTeadDown() {
	// deliberately does nothing
}

func (runner *ScenarioRunner) Concurrently(concurrently bool) *ScenarioRunner {
	runner.concurrently = concurrently
	return runner
}

func (runner *ScenarioRunner) Run() error {
	runner.logScenarioStartMessage()
	runner.startTime = Now()

	var runError error
	if runner.concurrently {
		runError = runner.runConcurrently()
	} else {
		runError = runner.runSequentially()
	}

	runner.finishTime = Now()
	runner.logHandler.Info("Finished running scenario \"" + runner.name + "\"")

	runner.logHandler.Info(runner.generateElapsedTimeString())

	runner.tearDown()

	return runError
}

func (runner *ScenarioRunner) LogHandler() handlers.LogHandler {
	return runner.logHandler
}

func (runner *ScenarioRunner) logScenarioStartMessage() {
	var runTypeText string
	if runner.concurrently {
		runTypeText = "concurrently"
	} else {
		runTypeText = "sequentially"
	}

	message := fmt.Sprintf("Scenario \"%s\": configured for %d run(s), executed %s", runner.name, runner.runNumber, runTypeText)
	runner.logHandler.Info(message)
}

func (runner *ScenarioRunner) generateElapsedTimeString() string {
	return fmt.Sprintf("Total elapsed time of scenario = [%v]", runner.ElapsedTime())
}

func (runner *ScenarioRunner) ElapsedTime() Duration {
	return runner.finishTime.Sub(runner.startTime)
}

func (runner *ScenarioRunner) runConcurrently() error {
	var wg sync.WaitGroup

	runThenDone := func(runNumber uint64) {
		runner.run(runNumber)
		wg.Done()
	}

	wg.Add(int(runner.runNumber))
	for runNumber := uint64(1); runNumber <= runner.runNumber; runNumber++ {
		go runThenDone(runNumber)
	}
	wg.Wait()

	return nil
}

func (runner *ScenarioRunner) runSequentially() error {
	for runNumber := uint64(1); runNumber <= runner.runNumber; runNumber++ {
		runner.run(runNumber)
	}
	return nil
}

func (runner *ScenarioRunner) run(runNumber uint64) {
	annealerCopy := runner.annealer.Clone()

	annealerCopy.SetId(runner.generateCloneId(runNumber))

	runner.logRunStartMessage(runNumber)
	annealerCopy.Anneal()
	runner.logRunFinishedMessage(runNumber)
}

func (runner *ScenarioRunner) generateCloneId(runNumber uint64) string {
	if runner.runNumber > 1 {
		return fmt.Sprintf("%s (%d/%d)", runner.name, runNumber, runner.runNumber)
	} else {
		return runner.name
	}
}

func (runner *ScenarioRunner) logRunStartMessage(runNumber uint64) {
	if runner.runNumber > 1 {
		runner.logHandler.Info(runner.generateCloneId(runNumber) + ": run started")
	}
}

func (runner *ScenarioRunner) logRunFinishedMessage(runNumber uint64) {
	if runner.runNumber > 1 {
		runner.logHandler.Info(runner.generateCloneId(runNumber) + ": run finished")
	}
}

type ProfilableScenarioRunner struct {
	base        CallableScenarioRunner
	profilePath string
}

func (runner *ProfilableScenarioRunner) ThatProfiles(base CallableScenarioRunner) *ProfilableScenarioRunner {
	runner.base = base
	return runner
}

func (runner *ProfilableScenarioRunner) ToFile(filePath string) *ProfilableScenarioRunner {
	runner.profilePath = filePath
	return runner
}

func (runner *ProfilableScenarioRunner) LogHandler() handlers.LogHandler {
	return runner.base.LogHandler()
}

func (runner *ProfilableScenarioRunner) Run() error {
	runner.LogHandler().Info("About to collect cpu profiling data to file [" + runner.profilePath + "]")
	defer runner.LogHandler().Info("Collection of cpu profiling data to file [" + runner.profilePath + "] complete.")

	return profiling.CpuProfileOfFunctionToFile(runner.base.Run, runner.profilePath)
}

type OleSafeScenarioRunner struct {
	base CallableScenarioRunner
}

func (runner *OleSafeScenarioRunner) ThatLocks(base CallableScenarioRunner) *OleSafeScenarioRunner {
	runner.base = base
	return runner
}

func (runner *OleSafeScenarioRunner) LogHandler() handlers.LogHandler {
	return runner.base.LogHandler()
}

func (runner *OleSafeScenarioRunner) Run() error {
	runner.LogHandler().Debug("Making scenario runner goroutine OLE thread-safe")
	runtime.LockOSThread()
	ole.CoInitialize(0)

	defer ole.CoUninitialize()
	defer runtime.UnlockOSThread()
	defer runner.LogHandler().Debug("Released OLE thread-safe scenario runner goroutine resources")

	return runner.base.Run()
}
