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

	startTime  Time
	finishTime Time
}

func (ao *ScenarioRunner) ForAnnealer(annealer shared.Annealer) *ScenarioRunner {
	ao.runNumber = 1
	ao.concurrently = false
	ao.name = "Default Scenario"

	ao.logHandler = annealer.LogHandler()
	ao.annealer = annealer

	return ao
}

func (ao *ScenarioRunner) WithName(name string) *ScenarioRunner {
	if name == "" {
		return ao
	}
	ao.name = name
	return ao
}

func (ao *ScenarioRunner) WithRunNumber(runNumber uint64) *ScenarioRunner {
	if runNumber == 0 {
		return ao
	}
	ao.runNumber = runNumber
	return ao
}

func (ao *ScenarioRunner) Concurrently(concurrently bool) *ScenarioRunner {
	ao.concurrently = concurrently
	return ao
}

func (ao *ScenarioRunner) Run() error {
	ao.logScenarioStartMessage()
	ao.startTime = Now()

	var runError error
	if ao.concurrently {
		runError = ao.runConcurrently()
	} else {
		runError = ao.runSequentially()
	}

	ao.finishTime = Now()
	ao.logHandler.Info("Finished running scenario \"" + ao.name + "\"")

	ao.logHandler.Info(ao.generateElapsedTimeString())

	return runError
}

func (ao *ScenarioRunner) LogHandler() handlers.LogHandler {
	return ao.logHandler
}

func (ao *ScenarioRunner) logScenarioStartMessage() {
	var runTypeText string
	if ao.concurrently {
		runTypeText = "concurrently"
	} else {
		runTypeText = "sequentially"
	}

	message := fmt.Sprintf("Scenario \"%s\": configured for %d run(s), executed %s", ao.name, ao.runNumber, runTypeText)
	ao.logHandler.Info(message)
}

func (ao *ScenarioRunner) generateElapsedTimeString() string {
	return fmt.Sprintf("Total elapsed time of scenario = [%v]", ao.ElapsedTime())
}

func (ao *ScenarioRunner) ElapsedTime() Duration {
	return ao.finishTime.Sub(ao.startTime)
}

func (ao *ScenarioRunner) runConcurrently() error {
	var wg sync.WaitGroup

	runThenDone := func(runNumber uint64) {
		ao.run(runNumber)
		wg.Done()
	}

	wg.Add(int(ao.runNumber))
	for runNumber := uint64(1); runNumber <= ao.runNumber; runNumber++ {
		go runThenDone(runNumber)
	}
	wg.Wait()

	return nil
}

func (ao *ScenarioRunner) runSequentially() error {
	for runNumber := uint64(1); runNumber <= ao.runNumber; runNumber++ {
		ao.run(runNumber)
	}
	return nil
}

func (ao *ScenarioRunner) run(runNumber uint64) {
	annealerCopy := ao.annealer.Clone()

	annealerCopy.SetId(ao.generateCloneId(runNumber))

	ao.logRunStartMessage(runNumber)
	annealerCopy.Anneal()
	ao.logRunFinishedMessage(runNumber)
}

func (ao *ScenarioRunner) generateCloneId(runNumber uint64) string {
	if ao.runNumber > 1 {
		return fmt.Sprintf("%s (%d/%d)", ao.name, runNumber, ao.runNumber)
	} else {
		return ao.name
	}
}

func (ao *ScenarioRunner) logRunStartMessage(runNumber uint64) {
	if ao.runNumber > 1 {
		ao.logHandler.Info(ao.generateCloneId(runNumber) + ": run started")
	}
}

func (ao *ScenarioRunner) logRunFinishedMessage(runNumber uint64) {
	if ao.runNumber > 1 {
		ao.logHandler.Info(ao.generateCloneId(runNumber) + ": run finished")
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
	runner.LogHandler().Debug("Locking scenario runner goroutine to the OS thread and co-initialising OLE")
	runtime.LockOSThread()
	ole.CoInitialize(0)

	defer ole.CoUninitialize()
	defer runtime.UnlockOSThread()
	defer runner.LogHandler().Debug("Unlocked scenario runner goroutine from the OS thread and Co-unitialised OLE")

	return runner.base.Run()
}
