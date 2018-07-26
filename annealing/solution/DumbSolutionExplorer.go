// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

type DumbSolutionExplorer struct {
	BaseSolutionExplorer
}

func (this *DumbSolutionExplorer) SetObjectiveValue(initialObjectiveValue float64) {
	this.objectiveValue = initialObjectiveValue
}

func (this *DumbSolutionExplorer) Initialise() {
	this.BaseSolutionExplorer.Initialise()
}

func (this *DumbSolutionExplorer) TearDown() {
	this.BaseSolutionExplorer.TearDown()
}

func (this *DumbSolutionExplorer) TryRandomChange(temperature float64) {
	this.makeRandomChange()
	DecideOnWhetherToAcceptChange(this, temperature)
}

func (this *DumbSolutionExplorer) makeRandomChange() {
	randomValue := this.RandomNumberGenerator().Intn(2)

	var changeInObjectiveValue float64
	switch randomValue {
	case 0:
		changeInObjectiveValue = -1
	case 1:
		changeInObjectiveValue = 1
	}
	if this.objectiveValue+changeInObjectiveValue >= 0 {
		this.changeInObjectiveValue = changeInObjectiveValue
	} else {
		this.changeInObjectiveValue = 0
	}

	this.objectiveValue += this.changeInObjectiveValue
}

func (this *DumbSolutionExplorer) AcceptLastChange() {
	this.BaseSolutionExplorer.AcceptLastChange()
}

func (this *DumbSolutionExplorer) RevertLastChange() {
	this.objectiveValue -= this.changeInObjectiveValue
	this.BaseSolutionExplorer.RevertLastChange()
}
