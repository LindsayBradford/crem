// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import 	. "github.com/LindsayBradford/crm/annealing/objectives"

type KnapsackObjectiveManager struct {
	BaseObjectiveManager
}

func (this *KnapsackObjectiveManager) Initialise() {
	this.BaseObjectiveManager.Initialise()
}

func (this *KnapsackObjectiveManager) TryRandomChange(temperature float64) {
	this.makeRandomChange()
	DecideOnWhetherToAcceptChange(this, temperature)
}

func (this *KnapsackObjectiveManager) makeRandomChange() {
	randomValue := this.RandomNumberGenerator().Intn(2)

	var changeInObjectiveValue float64
	switch randomValue {
	case 0:
		changeInObjectiveValue = -1
	case 1:
		changeInObjectiveValue = 1
	}

	if this.ObjectiveValue() + changeInObjectiveValue >= 0 {
		this.SetChangeInObjectiveValue(changeInObjectiveValue)
	} else {
		this.SetChangeInObjectiveValue(0)
	}
	this.SetObjectiveValue(this.ObjectiveValue() + this.ChangeInObjectiveValue())
}

func (this *KnapsackObjectiveManager) AcceptLastChange()  {
	this.BaseObjectiveManager.AcceptLastChange()
}

func (this *KnapsackObjectiveManager) RevertLastChange()  {
	this.SetObjectiveValue(this.ObjectiveValue() - this.ChangeInObjectiveValue())
	this.BaseObjectiveManager.RevertLastChange()
}

