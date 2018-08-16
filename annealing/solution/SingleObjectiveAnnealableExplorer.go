// Copyright (c) 2018 Australian Rivers Institute.

package solution

import "math"

type SingleObjectiveAnnealableExplorer struct {
	BaseExplorer
}

func (explorer *SingleObjectiveAnnealableExplorer) TryRandomChange(temperature float64) {
	explorer.makeRandomChange()
	explorer.DecideOnWhetherToAcceptChange(temperature)
}

func (explorer *SingleObjectiveAnnealableExplorer) makeRandomChange() {}

func (explorer *SingleObjectiveAnnealableExplorer) DecideOnWhetherToAcceptChange(annealingTemperature float64) {
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
