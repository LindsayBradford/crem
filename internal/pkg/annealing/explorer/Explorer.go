// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package explorer

import (
	"math/rand"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/pkg/logging"
)

type Explorer interface {
	Name() string

	ScenarioId() string
	SetScenarioId(name string)

	Model() model.Model
	SetModel(model model.Model)

	Initialise()
	Clone() Explorer
	TryRandomChange(temperature float64)

	ObjectiveValue() float64
	SetObjectiveValue(objectiveValue float64)

	ChangeInObjectiveValue() float64
	SetChangeInObjectiveValue(change float64)

	AcceptanceProbability() float64
	SetAcceptanceProbability(probability float64)

	ChangeIsDesirable() bool
	DecideOnWhetherToAcceptChange(annealingTemperature float64, acceptFunction func(), rejectFunction func())
	AcceptLastChange()
	ChangeAccepted() bool
	RevertLastChange()

	SetRandomNumberGenerator(*rand.Rand)
	RandomNumberGenerator() *rand.Rand

	SetLogHandler(logger logging.Logger) error
	LogHandler() logging.Logger

	TearDown()
}

type AnnealableExplorer interface {
	ObjectiveValue() float64

	ChangeInObjectiveValue() float64

	AcceptanceProbability() float64

	ChangeIsDesirable() bool
	ChangeAccepted() bool
}
