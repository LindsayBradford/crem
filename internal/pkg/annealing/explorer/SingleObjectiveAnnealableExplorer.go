// Copyright (c) 2018 Australian Rivers Institute.

package explorer

import (
	"math"
)

type SingleObjectiveAnnealableExplorer struct {
	BaseExplorer
}

func (explorer *SingleObjectiveAnnealableExplorer) TryRandomChange(temperature float64) {
	explorer.makeRandomChange()
	explorer.DecideOnWhetherToAcceptChange(temperature, explorer.AcceptLastChange, explorer.RevertLastChange)
}

func (explorer *SingleObjectiveAnnealableExplorer) makeRandomChange() {}

func (explorer *SingleObjectiveAnnealableExplorer) DecideOnWhetherToAcceptChange(annealingTemperature float64, acceptChange func(), revertChange func()) {
	if explorer.ChangeIsDesirable() {
		explorer.SetAcceptanceProbability(1)
		acceptChange()
	} else {
		probabilityToAcceptBadChange := math.Exp(-explorer.ChangeInObjectiveValue() / annealingTemperature)
		explorer.SetAcceptanceProbability(probabilityToAcceptBadChange)

		randomValue := newRandomValue(explorer.RandomNumberGenerator())
		if probabilityToAcceptBadChange > randomValue {
			acceptChange()
		} else {
			revertChange()
		}
	}
}
