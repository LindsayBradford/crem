// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

type DumbSolutionTourer struct {
	BaseSolutionTourer
}

func (this *DumbSolutionTourer) SetObjectiveValue(initialObjectiveValue float64) {
	this.objectiveValue = initialObjectiveValue
}

func (this *DumbSolutionTourer) Initialise() {
	this.BaseSolutionTourer.Initialise()
}

func (this *DumbSolutionTourer) TearDown() {
	this.BaseSolutionTourer.TearDown()
}

func (this *DumbSolutionTourer) TryRandomChange(temperature float64) {
	this.makeRandomChange()
	DecideOnWhetherToAcceptChange(this, temperature)
}

func (this *DumbSolutionTourer) makeRandomChange() {
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

func (this *DumbSolutionTourer) AcceptLastChange()  {
	this.BaseSolutionTourer.AcceptLastChange()
}

func (this *DumbSolutionTourer) RevertLastChange()  {
	this.objectiveValue -= this.changeInObjectiveValue
	this.BaseSolutionTourer.RevertLastChange()
}
