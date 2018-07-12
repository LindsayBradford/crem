// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package objectives

import (
	"math"
	"math/rand"
	"time"
)

type ObjectiveManager interface {
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
}

type BaseObjectiveManager struct {
	objectiveValue         float64
	changeInObjectiveValue float64
	changeIsDesirable      bool
	changeAccepted         bool
	acceptanceProbability  float64
	randomNumberGenerator  *rand.Rand
}

func (this *BaseObjectiveManager) Initialise() {
	this.SetRandomNumberGenerator(rand.New(rand.NewSource(time.Now().UnixNano())))
}

func (this *BaseObjectiveManager) RandomNumberGenerator() *rand.Rand {
	return this.randomNumberGenerator
}

func (this *BaseObjectiveManager) SetRandomNumberGenerator(generator *rand.Rand) {
	this.randomNumberGenerator = generator
}

func (this *BaseObjectiveManager) TryRandomChange(temperature float64) {
	this.makeRandomChange()
	DecideOnWhetherToAcceptChange(this, temperature)
}

func (this *BaseObjectiveManager) SetObjectiveValue(objectiveValue float64) {
	this.objectiveValue = objectiveValue
}

func (this *BaseObjectiveManager) ObjectiveValue() float64 {
	return this.objectiveValue
}

func (this *BaseObjectiveManager) ChangeInObjectiveValue() float64 {
	return this.changeInObjectiveValue
}

func (this *BaseObjectiveManager) SetChangeInObjectiveValue(change float64) {
	this.changeInObjectiveValue = change
}

func (this *BaseObjectiveManager) AcceptanceProbability() float64 {
	return this.acceptanceProbability
}

func (this *BaseObjectiveManager) SetAcceptanceProbability(probability float64) {
	this.acceptanceProbability = probability
}

func (this *BaseObjectiveManager) makeRandomChange() {}

func DecideOnWhetherToAcceptChange(manager ObjectiveManager,  annealingTemperature float64) {
	if (manager.ChangeIsDesirable()) {
		manager.SetAcceptanceProbability(1)
		manager.AcceptLastChange()
	} else {
		probabilityToAcceptBadChange := math.Exp(-manager.ChangeInObjectiveValue() / annealingTemperature)
		manager.SetAcceptanceProbability(probabilityToAcceptBadChange)

		randomValue  := newRandomValue(manager.RandomNumberGenerator())
		if probabilityToAcceptBadChange > randomValue {
			manager.AcceptLastChange()
		} else {
			manager.RevertLastChange()
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

func (this *BaseObjectiveManager) ChangeIsDesirable() bool {
	if this.changeInObjectiveValue <= 0  {
		return true
	}
	return false
}

func (this *BaseObjectiveManager) AcceptLastChange()  {
	this.changeAccepted = true
}

func (this *BaseObjectiveManager) RevertLastChange()  {
	this.changeAccepted = false
}

func (this *BaseObjectiveManager) ChangeAccepted() bool {
	return this.changeAccepted
}