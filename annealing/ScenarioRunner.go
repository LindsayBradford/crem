// Copyright (c) 2018 Australian Rivers Institute.

package annealing

import (
	"fmt"
	"sync"
	. "time"

	"github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/logging/handlers"
)

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

func (ao *ScenarioRunner) Run() {
	ao.logScenarioStartMessage()
	ao.startTime = Now()

	if ao.concurrently {
		ao.runConcurrently()
	} else {
		ao.runSequentially()
	}

	ao.finishTime = Now()
	ao.logHandler.Info("Finished running scenario \"" + ao.name + "\"")

	ao.logHandler.Info(ao.generateElapsedTimeString())
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

func (ao *ScenarioRunner) runConcurrently() {
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
}

func (ao *ScenarioRunner) runSequentially() {
	for runNumber := uint64(1); runNumber <= ao.runNumber; runNumber++ {
		ao.run(runNumber)
	}
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
