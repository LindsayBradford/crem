// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

import (
	"errors"
	"math"
	"math/rand"
	"time"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

type SolutionTourer interface {
	Initialise()
	TryRandomChange(temperature float64)

	ObjectiveValue() float64
	SetObjectiveValue(objectiveValue float64)

	ChangeInObjectiveValue() float64
	SetChangeInObjectiveValue(change float64)

	AcceptanceProbability() float64
	SetAcceptanceProbability(probability float64)

	ChangeIsDesirable() bool
	AcceptLastChange()
	ChangeAccepted() bool
	RevertLastChange()

	SetRandomNumberGenerator(*rand.Rand)
	RandomNumberGenerator() *rand.Rand

	SetLogHandler(logger LogHandler) error
	LogHandler() LogHandler

	TearDown()
}

type BaseSolutionTourer struct {
	objectiveValue         float64
	changeInObjectiveValue float64
	changeIsDesirable      bool
	changeAccepted         bool
	acceptanceProbability  float64
	randomNumberGenerator  *rand.Rand
	logHandler             LogHandler
}

func (this *BaseSolutionTourer) Initialise() {
	this.logHandler.Debug("Initialising Solution Tourer")
	this.SetRandomNumberGenerator(rand.New(rand.NewSource(time.Now().UnixNano())))
}

func (this *BaseSolutionTourer) TearDown() {
	this.logHandler.Debug("Triggering tear-down of Solution Tourer")
}

func (this *BaseSolutionTourer) RandomNumberGenerator() *rand.Rand {
	return this.randomNumberGenerator
}

func (this *BaseSolutionTourer) SetRandomNumberGenerator(generator *rand.Rand) {
	this.randomNumberGenerator = generator
}

func (this *BaseSolutionTourer) LogHandler() LogHandler {
	return this.logHandler
}

func (this *BaseSolutionTourer) SetLogHandler(logHandler LogHandler) error {
	if logHandler == nil {
		return errors.New("Invalid attempt to set log handler to nil value")
	}
	this.logHandler = logHandler
	return nil
}

func (this *BaseSolutionTourer) TryRandomChange(temperature float64) {
	this.makeRandomChange()
	DecideOnWhetherToAcceptChange(this, temperature)
}

func (this *BaseSolutionTourer) SetObjectiveValue(objectiveValue float64) {
	this.objectiveValue = objectiveValue
}

func (this *BaseSolutionTourer) ObjectiveValue() float64 {
	return this.objectiveValue
}

func (this *BaseSolutionTourer) ChangeInObjectiveValue() float64 {
	return this.changeInObjectiveValue
}

func (this *BaseSolutionTourer) SetChangeInObjectiveValue(change float64) {
	this.changeInObjectiveValue = change
}

func (this *BaseSolutionTourer) AcceptanceProbability() float64 {
	return this.acceptanceProbability
}

func (this *BaseSolutionTourer) SetAcceptanceProbability(probability float64) {
	this.acceptanceProbability = probability
}

func (this *BaseSolutionTourer) makeRandomChange() {}

func DecideOnWhetherToAcceptChange(tourer SolutionTourer,  annealingTemperature float64) {
	if (tourer.ChangeIsDesirable()) {
		tourer.SetAcceptanceProbability(1)
		tourer.AcceptLastChange()
	} else {
		probabilityToAcceptBadChange := math.Exp(-tourer.ChangeInObjectiveValue() / annealingTemperature)
		tourer.SetAcceptanceProbability(probabilityToAcceptBadChange)

		randomValue  := newRandomValue(tourer.RandomNumberGenerator())
		if probabilityToAcceptBadChange > randomValue {
			tourer.AcceptLastChange()
		} else {
			tourer.RevertLastChange()
		}
	}
}

// newRandomValue returns the next random number in the range [0,1] from the supplied randomNumberGenerator.
// (which by default returns a random number in the range [0,1).
// See: http://mumble.net/~campbell/2014/04/28/uniform-random-float
func newRandomValue(randomNumberGenerator *rand.Rand) float64 {
	distributionRange := int64(math.Pow(2,53))
	return float64(randomNumberGenerator.Int63n(distributionRange)) / float64(distributionRange - 1)
}

func (this *BaseSolutionTourer) ChangeIsDesirable() bool {
	if this.changeInObjectiveValue <= 0  {
		return true
	}
	return false
}

func (this *BaseSolutionTourer) AcceptLastChange()  {
	this.changeAccepted = true
}

func (this *BaseSolutionTourer) RevertLastChange()  {
	this.changeAccepted = false
}

func (this *BaseSolutionTourer) ChangeAccepted() bool {
	return this.changeAccepted
}