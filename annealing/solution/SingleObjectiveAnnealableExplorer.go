// Copyright (c) 2018 Australian Rivers Institute.

package solution

type SingleObjectiveAnnealableExplorer struct {
	BaseExplorer
}

func (explorer *SingleObjectiveAnnealableExplorer) TryRandomChange(temperature float64) {
	explorer.makeRandomChange()
	explorer.DecideOnWhetherToAcceptChange(temperature)
}

func (explorer *SingleObjectiveAnnealableExplorer) makeRandomChange() {}
