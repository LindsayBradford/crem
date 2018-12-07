// Copyright (c) 2018 Australian Rivers Institute.

package explorer

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

type SingleObjectiveAnnealableExplorer struct {
	name.ContainedName
	scenarioId string
	model.ContainedModel
	logging.ContainedLogger

	objectiveValue         float64
	changeInObjectiveValue float64
	changeIsDesirable      bool
	changeAccepted         bool
	acceptanceProbability  float64
	rand.ContainedRand
}

func (explorer *SingleObjectiveAnnealableExplorer) TryRandomChange(temperature float64) {
	explorer.makeRandomChange()
	explorer.DecideOnWhetherToAcceptChange(temperature, explorer.AcceptLastChange, explorer.RevertLastChange)
}

func (explorer *SingleObjectiveAnnealableExplorer) makeRandomChange() {}

func (explorer *SingleObjectiveAnnealableExplorer) DecideOnWhetherToAcceptChange(annealingTemperature float64, acceptChange func(), revertChange func()) {
	if explorer.ChangeIsDesirable() {
		explorer.SetAcceptanceProbability(1)
		acceptChange()
	} else {
		probabilityToAcceptBadChange := math.Exp(-explorer.ChangeInObjectiveValue() / annealingTemperature)
		explorer.SetAcceptanceProbability(probabilityToAcceptBadChange)

		randomValue := explorer.RandomNumberGenerator().Float64Unitary()
		if probabilityToAcceptBadChange > randomValue {
			acceptChange()
		} else {
			revertChange()
		}
	}
}

func (explorer *SingleObjectiveAnnealableExplorer) Initialise() {
	explorer.LogHandler().Debug(explorer.scenarioId + ": Initialising Solution Explorer")
	explorer.SetRandomNumberGenerator(rand.NewTimeSeeded())
}

func (explorer *SingleObjectiveAnnealableExplorer) TearDown() {
	explorer.LogHandler().Debug(explorer.scenarioId + ": Triggering tear-down of Solution Explorer")
}

func (explorer *SingleObjectiveAnnealableExplorer) WithName(name string) *SingleObjectiveAnnealableExplorer {
	explorer.SetName(name)
	return explorer
}

func (explorer *SingleObjectiveAnnealableExplorer) WithModel(model model.Model) *SingleObjectiveAnnealableExplorer {
	explorer.SetModel(model)
	return explorer
}

func (explorer *SingleObjectiveAnnealableExplorer) ScenarioId() string {
	return explorer.scenarioId
}

func (explorer *SingleObjectiveAnnealableExplorer) SetScenarioId(id string) {
	explorer.scenarioId = id
}

func (explorer *SingleObjectiveAnnealableExplorer) WithScenarioId(id string) *SingleObjectiveAnnealableExplorer {
	explorer.scenarioId = id
	return explorer
}

func (explorer *SingleObjectiveAnnealableExplorer) SetObjectiveValue(objectiveValue float64) {
	explorer.objectiveValue = objectiveValue
}

func (explorer *SingleObjectiveAnnealableExplorer) ObjectiveValue() float64 {
	return explorer.objectiveValue
}

func (explorer *SingleObjectiveAnnealableExplorer) ChangeInObjectiveValue() float64 {
	return explorer.changeInObjectiveValue
}

func (explorer *SingleObjectiveAnnealableExplorer) SetChangeInObjectiveValue(change float64) {
	explorer.changeInObjectiveValue = change
}

func (explorer *SingleObjectiveAnnealableExplorer) AcceptanceProbability() float64 {
	return explorer.acceptanceProbability
}

func (explorer *SingleObjectiveAnnealableExplorer) SetAcceptanceProbability(probability float64) {
	explorer.acceptanceProbability = probability
}

func (explorer *SingleObjectiveAnnealableExplorer) ChangeIsDesirable() bool {
	if explorer.changeInObjectiveValue <= 0 {
		return true
	}
	return false
}

func (explorer *SingleObjectiveAnnealableExplorer) AcceptLastChange() {
	explorer.changeAccepted = true
}

func (explorer *SingleObjectiveAnnealableExplorer) RevertLastChange() {
	explorer.changeAccepted = false
}

func (explorer *SingleObjectiveAnnealableExplorer) ChangeAccepted() bool {
	return explorer.changeAccepted
}

func (explorer *SingleObjectiveAnnealableExplorer) DeepClone() Explorer {
	clone := *explorer
	modelClone := clone.Model().DeepClone()
	clone.SetModel(modelClone)
	return &clone
}

func (explorer *SingleObjectiveAnnealableExplorer) CloneObservable() Explorer {
	clone := *explorer
	clone.SetModel(model.NullModel)
	return &clone
}
