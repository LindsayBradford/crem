// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

import (
	"errors"
	. "github.com/LindsayBradford/crm/logging/handlers"
	"math"
	"math/rand"
	"time"
)

type SolutionExplorer interface {
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

type BaseSolutionExplorer struct {
	objectiveValue         float64
	changeInObjectiveValue float64
	changeIsDesirable      bool
	changeAccepted         bool
	acceptanceProbability  float64
	randomNumberGenerator  *rand.Rand
	logHandler             LogHandler
}

func (this *BaseSolutionExplorer) Initialise() {
	this.logHandler.Debug("Initialising Solution Explorer")
	this.SetRandomNumberGenerator(rand.New(rand.NewSource(time.Now().UnixNano())))
}

func (this *BaseSolutionExplorer) TearDown() {
	this.logHandler.Debug("Triggering tear-down of Solution Explorer")
}

func (this *BaseSolutionExplorer) RandomNumberGenerator() *rand.Rand {
	return this.randomNumberGenerator
}

func (this *BaseSolutionExplorer) SetRandomNumberGenerator(generator *rand.Rand) {
	this.randomNumberGenerator = generator
}

func (this *BaseSolutionExplorer) LogHandler() LogHandler {
	return this.logHandler
}

func (this *BaseSolutionExplorer) SetLogHandler(logHandler LogHandler) error {
	if logHandler == nil {
		return errors.New("Invalid attempt to set log handler to nil value")
	}
	this.logHandler = logHandler
	return nil
}

func (this *BaseSolutionExplorer) TryRandomChange(temperature float64) {
	this.makeRandomChange()
	DecideOnWhetherToAcceptChange(this, temperature)
}

func (this *BaseSolutionExplorer) SetObjectiveValue(objectiveValue float64) {
	this.objectiveValue = objectiveValue
}

func (this *BaseSolutionExplorer) ObjectiveValue() float64 {
	return this.objectiveValue
}

func (this *BaseSolutionExplorer) ChangeInObjectiveValue() float64 {
	return this.changeInObjectiveValue
}

func (this *BaseSolutionExplorer) SetChangeInObjectiveValue(change float64) {
	this.changeInObjectiveValue = change
}

func (this *BaseSolutionExplorer) AcceptanceProbability() float64 {
	return this.acceptanceProbability
}

func (this *BaseSolutionExplorer) SetAcceptanceProbability(probability float64) {
	this.acceptanceProbability = probability
}

func (this *BaseSolutionExplorer) makeRandomChange() {}

func DecideOnWhetherToAcceptChange(explorer SolutionExplorer, annealingTemperature float64) {
	if explorer.ChangeIsDesirable() {
		explorer.SetAcceptanceProbability(1)
		explorer.AcceptLastChange()
	} else {
		probabilityToAcceptBadChange := math.Exp(-explorer.ChangeInObjectiveValue() / annealingTemperature)
		explorer.SetAcceptanceProbability(probabilityToAcceptBadChange)

		randomValue := newRandomValue(explorer.RandomNumberGenerator())
		if probabilityToAcceptBadChange > randomValue {
			explorer.AcceptLastChange()
		} else {
			explorer.RevertLastChange()
		}
	}
}

// newRandomValue returns the next random number in the range [0,1] from the supplied randomNumberGenerator.
// (which by default returns a random number in the range [0,1).
// See: http://mumble.net/~campbell/2014/04/28/uniform-random-float
func newRandomValue(randomNumberGenerator *rand.Rand) float64 {
	distributionRange := int64(math.Pow(2, 53))
	return float64(randomNumberGenerator.Int63n(distributionRange)) / float64(distributionRange-1)
}

func (this *BaseSolutionExplorer) ChangeIsDesirable() bool {
	if this.changeInObjectiveValue <= 0 {
		return true
	}
	return false
}

func (this *BaseSolutionExplorer) AcceptLastChange() {
	this.changeAccepted = true
}

func (this *BaseSolutionExplorer) RevertLastChange() {
	this.changeAccepted = false
}

func (this *BaseSolutionExplorer) ChangeAccepted() bool {
	return this.changeAccepted
}
