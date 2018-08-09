// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

type DumbSolutionExplorer struct {
	BaseSolutionExplorer
}

func (dse *DumbSolutionExplorer) WithName(name string) *DumbSolutionExplorer {
	dse.BaseSolutionExplorer.WithName(name)
	return dse
}

func (dse *DumbSolutionExplorer) SetObjectiveValue(initialObjectiveValue float64) {
	dse.objectiveValue = initialObjectiveValue
}

func (dse *DumbSolutionExplorer) Initialise() {
	dse.BaseSolutionExplorer.Initialise()
}

func (dse *DumbSolutionExplorer) TearDown() {
	dse.BaseSolutionExplorer.TearDown()
}

func (dse *DumbSolutionExplorer) TryRandomChange(temperature float64) {
	dse.makeRandomChange()
	DecideOnWhetherToAcceptChange(dse, temperature)
}

func (dse *DumbSolutionExplorer) makeRandomChange() {
	randomValue := dse.RandomNumberGenerator().Intn(2)

	var changeInObjectiveValue float64
	switch randomValue {
	case 0:
		changeInObjectiveValue = -1
	case 1:
		changeInObjectiveValue = 1
	}
	if dse.objectiveValue+changeInObjectiveValue >= 0 {
		dse.changeInObjectiveValue = changeInObjectiveValue
	} else {
		dse.changeInObjectiveValue = 0
	}

	dse.objectiveValue += dse.changeInObjectiveValue
}

func (dse *DumbSolutionExplorer) AcceptLastChange() {
	dse.BaseSolutionExplorer.AcceptLastChange()
}

func (dse *DumbSolutionExplorer) RevertLastChange() {
	dse.objectiveValue -= dse.changeInObjectiveValue
	dse.BaseSolutionExplorer.RevertLastChange()
}
