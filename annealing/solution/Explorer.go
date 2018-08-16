// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

import (
	. "github.com/LindsayBradford/crm/logging/handlers"
	"math/rand"
)

type Explorer interface {
	Name() string
	SetName(name string)

	Initialise()
	TryRandomChange(temperature float64)

	ObjectiveValue() float64
	SetObjectiveValue(objectiveValue float64)

	ChangeInObjectiveValue() float64
	SetChangeInObjectiveValue(change float64)

	AcceptanceProbability() float64
	SetAcceptanceProbability(probability float64)

	ChangeIsDesirable() bool
	DecideOnWhetherToAcceptChange(annealingTemperature float64)
	AcceptLastChange()
	ChangeAccepted() bool
	RevertLastChange()

	SetRandomNumberGenerator(*rand.Rand)
	RandomNumberGenerator() *rand.Rand

	SetLogHandler(logger LogHandler) error
	LogHandler() LogHandler

	TearDown()
}

type AnnealableExplorer interface {
	ObjectiveValue() float64

	ChangeInObjectiveValue() float64

	AcceptanceProbability() float64

	ChangeIsDesirable() bool
	ChangeAccepted() bool
}
