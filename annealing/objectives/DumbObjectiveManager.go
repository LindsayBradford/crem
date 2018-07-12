// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package objectives

type DumbObjectiveManager struct {
	BaseObjectiveManager
}

func (this *DumbObjectiveManager) SetObjectiveValue(initialObjectiveValue float64) {
	this.objectiveValue = initialObjectiveValue
}

func (this *DumbObjectiveManager) Initialise() {
	this.BaseObjectiveManager.Initialise()
}

func (this *DumbObjectiveManager) TryRandomChange(temperature float64) {
	this.makeRandomChange()
	DecideOnWhetherToAcceptChange(this, temperature)
}

func (this *DumbObjectiveManager) makeRandomChange() {
	randomValue := this.RandomNumberGenerator().Intn(2)

	var changeInObjectiveValue float64
	switch randomValue {
	case 0:
		changeInObjectiveValue = -1
	case 1:
		changeInObjectiveValue = 1
	}
	if this.objectiveValue + changeInObjectiveValue >= 0 {
		this.changeInObjectiveValue = changeInObjectiveValue
	} else {
		this.changeInObjectiveValue = 0
	}

	this.objectiveValue += this.changeInObjectiveValue
}

func (this *DumbObjectiveManager) AcceptLastChange()  {
	this.BaseObjectiveManager.AcceptLastChange()
}

func (this *DumbObjectiveManager) RevertLastChange()  {
	this.objectiveValue -= this.changeInObjectiveValue
	this.BaseObjectiveManager.RevertLastChange()
}
