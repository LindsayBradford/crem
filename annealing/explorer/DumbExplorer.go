// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package explorer

type DumbExplorer struct {
	SingleObjectiveAnnealableExplorer
}

func (dse *DumbExplorer) WithName(name string) *DumbExplorer {
	dse.SingleObjectiveAnnealableExplorer.WithName(name)
	return dse
}

func (dse *DumbExplorer) SetObjectiveValue(initialObjectiveValue float64) {
	dse.objectiveValue = initialObjectiveValue
}

func (dse *DumbExplorer) Initialise() {
	dse.SingleObjectiveAnnealableExplorer.Initialise()
}

func (dse *DumbExplorer) TearDown() {
	dse.SingleObjectiveAnnealableExplorer.TearDown()
}

func (dse *DumbExplorer) TryRandomChange(temperature float64) {
	dse.makeRandomChange()
	dse.DecideOnWhetherToAcceptChange(temperature, dse.AcceptLastChange, dse.RevertLastChange)
}

func (dse *DumbExplorer) makeRandomChange() {
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

func (dse *DumbExplorer) AcceptLastChange() {
	dse.SingleObjectiveAnnealableExplorer.AcceptLastChange()
}

func (dse *DumbExplorer) RevertLastChange() {
	dse.objectiveValue -= dse.changeInObjectiveValue
	dse.SingleObjectiveAnnealableExplorer.RevertLastChange()
}

func (dse *DumbExplorer) Clone() Explorer {
	clone := *dse
	return &clone
}
