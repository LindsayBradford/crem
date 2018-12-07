// Copyright (c) 2018 Australian Rivers Institute.

package explorer

import (
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/pkg/errors"
)

type BaseExplorer struct {
	name.Named
	scenarioId string
	model      model.Model

	objectiveValue         float64
	changeInObjectiveValue float64
	changeIsDesirable      bool
	changeAccepted         bool
	acceptanceProbability  float64
	rand.ContainedRand
	logHandler logging.Logger
}

func (explorer *BaseExplorer) Initialise() {
	explorer.logHandler.Debug(explorer.scenarioId + ": Initialising Solution Explorer")
	explorer.SetRandomNumberGenerator(rand.NewTimeSeeded())
}

func (explorer *BaseExplorer) TearDown() {
	explorer.logHandler.Debug(explorer.scenarioId + ": Triggering tear-down of Solution Explorer")
}

func (explorer *BaseExplorer) WithName(name string) *BaseExplorer {
	explorer.SetName(name)
	return explorer
}

func (explorer *BaseExplorer) Model() model.Model {
	return explorer.model
}

func (explorer *BaseExplorer) SetModel(model model.Model) {
	explorer.model = model
}

func (explorer *BaseExplorer) WithModel(model model.Model) *BaseExplorer {
	explorer.model = model
	return explorer
}

func (explorer *BaseExplorer) ScenarioId() string {
	return explorer.scenarioId
}

func (explorer *BaseExplorer) SetScenarioId(id string) {
	explorer.scenarioId = id
}

func (explorer *BaseExplorer) WithScenarioId(id string) *BaseExplorer {
	explorer.scenarioId = id
	return explorer
}

func (explorer *BaseExplorer) LogHandler() logging.Logger {
	return explorer.logHandler
}

func (explorer *BaseExplorer) SetLogHandler(logHandler logging.Logger) error {
	if logHandler == nil {
		return errors.New("invalid attempt to set log handler to nil value")
	}
	explorer.logHandler = logHandler
	return nil
}

func (explorer *BaseExplorer) TryRandomChange(temperature float64) {}

func (explorer *BaseExplorer) SetObjectiveValue(objectiveValue float64) {
	explorer.objectiveValue = objectiveValue
}

func (explorer *BaseExplorer) ObjectiveValue() float64 {
	return explorer.objectiveValue
}

func (explorer *BaseExplorer) ChangeInObjectiveValue() float64 {
	return explorer.changeInObjectiveValue
}

func (explorer *BaseExplorer) SetChangeInObjectiveValue(change float64) {
	explorer.changeInObjectiveValue = change
}

func (explorer *BaseExplorer) AcceptanceProbability() float64 {
	return explorer.acceptanceProbability
}

func (explorer *BaseExplorer) SetAcceptanceProbability(probability float64) {
	explorer.acceptanceProbability = probability
}

func (explorer *BaseExplorer) DecideOnWhetherToAcceptChange(annealingTemperature float64, acceptFunction func(), rejectFunction func()) {
}

func (explorer *BaseExplorer) ChangeIsDesirable() bool {
	if explorer.changeInObjectiveValue <= 0 {
		return true
	}
	return false
}

func (explorer *BaseExplorer) AcceptLastChange() {
	explorer.changeAccepted = true
}

func (explorer *BaseExplorer) RevertLastChange() {
	explorer.changeAccepted = false
}

func (explorer *BaseExplorer) ChangeAccepted() bool {
	return explorer.changeAccepted
}

func (explorer *BaseExplorer) DeepClone() Explorer {
	clone := *explorer
	modelClone := clone.Model().DeepClone()
	clone.SetModel(modelClone)
	return &clone
}

func (explorer *BaseExplorer) CloneObservable() Explorer {
	clone := *explorer
	clone.SetModel(model.NullModel)
	return &clone
}
