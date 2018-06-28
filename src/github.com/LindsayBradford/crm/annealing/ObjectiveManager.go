// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealing

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	randomNumberGenerator *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type ObjectiveManager interface {
	Initialise(annealer *annealerBase)
	TryRandomChange()

	ChangeInObjectiveValue() float64
	SetChangeInObjectiveValue(change float64)

	ChangeIsDesirable() bool
	AcceptLastChange()
	RevertLastChange()

	NotifyObserversWithObjectiveEvaluation(note string)
}

type BaseObjectiveManager struct {
	annealer *annealerBase
	changeInObjectiveValue  float64
}

func (this *BaseObjectiveManager) Initialise(annealer *annealerBase) {
	this.annealer = annealer
}

func (this *BaseObjectiveManager) TryRandomChange() {
	this.makeRandomChange()
	DecideOnWhetherToAcceptChange(this, this.annealer.temperature)
}

func (this *BaseObjectiveManager) ChangeInObjectiveValue() float64 {
	return this.changeInObjectiveValue
}

func (this *BaseObjectiveManager) SetChangeInObjectiveValue(change float64) {
	this.changeInObjectiveValue = change
}

func (this *BaseObjectiveManager) NotifyObserversWithObjectiveEvaluation(note string) {
	this.annealer.notifyObserversWithObjectiveEvaluation(note)
}

func (this *BaseObjectiveManager) makeRandomChange() {}

func DecideOnWhetherToAcceptChange(manager ObjectiveManager,  annealingTemperature float64) {
	if (manager.ChangeIsDesirable()) {
		manager.AcceptLastChange()
	} else {
		probabilityToAcceptBadChange := math.Exp(-manager.ChangeInObjectiveValue() / annealingTemperature)
		randomValue  := newRandomValue()
		manager.NotifyObserversWithObjectiveEvaluation(
			"Deciding on whether to accept bad change. " +"Temp = [" +fmt.Sprintf("%f", annealingTemperature) +
				"], Probability to accept = [" + fmt.Sprintf("%f", probabilityToAcceptBadChange) +
				"], randomValue = [" + fmt.Sprintf("%f", randomValue) + "]")
		if probabilityToAcceptBadChange > randomValue {
			manager.AcceptLastChange()
		} else {
			manager.RevertLastChange()
		}
	}
}

// newRandomValue returns the next random number in the range [0,1] from the embedded randomNumberGenerator.
// (which by default returns a random number in the range [0,1).
// See: http://mumble.net/~campbell/2014/04/28/uniform-random-float
func newRandomValue() float64 {
	distributionRange := int64(math.Pow(2,53))
	return float64(randomNumberGenerator.Int63n(distributionRange)) / float64(distributionRange - 1)
}

func (this *BaseObjectiveManager) ChangeIsDesirable() bool {
	if this.changeInObjectiveValue < 0  {
		this.NotifyObserversWithObjectiveEvaluation("Change in objective value = [" + fmt.Sprintf("%f", this.changeInObjectiveValue) + "]. Change is desirable")
		return true
	}
	this.NotifyObserversWithObjectiveEvaluation("Change in objective value = [" + fmt.Sprintf("%f", this.changeInObjectiveValue) + "]. Change is NOT desirable")
	return false
}

func (this *BaseObjectiveManager) AcceptLastChange()  {}

func (this *BaseObjectiveManager) RevertLastChange()  {}

type DumbObjectiveManager struct {
	BaseObjectiveManager
	currentObjectiveValue float64
}

func (this *DumbObjectiveManager) Initialise(annealer *annealerBase) {
	this.annealer = annealer
	this.NotifyObserversWithObjectiveEvaluation("Initialising ObjectiveManager")
	this.currentObjectiveValue = float64(50000)
	this.NotifyObserversWithObjectiveEvaluation("Current Objective Value: " + fmt.Sprintf("%f", this.currentObjectiveValue))
}

func (this *DumbObjectiveManager) TryRandomChange() {
	this.makeRandomChange()
	DecideOnWhetherToAcceptChange(this, this.annealer.temperature)
}

func (this *DumbObjectiveManager) makeRandomChange() {
	randomValue := randomNumberGenerator.Intn(2)
	switch randomValue {
	case 0:
		this.changeInObjectiveValue = -1
	case 1:
		this.changeInObjectiveValue = 1
	}
	this.currentObjectiveValue += this.changeInObjectiveValue
}

func (this *DumbObjectiveManager) AcceptLastChange()  {
	this.NotifyObserversWithObjectiveEvaluation("Accepting Change. Current Objective Value = [" + fmt.Sprintf("%f", this.currentObjectiveValue) + "]")
}

func (this *DumbObjectiveManager) RevertLastChange()  {
	this.currentObjectiveValue -= this.changeInObjectiveValue
	this.NotifyObserversWithObjectiveEvaluation("Reverting Change. Current Objective Value = [" + fmt.Sprintf("%f", this.currentObjectiveValue) + "]")
}