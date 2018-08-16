// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

var NULL_EXPLORER = new(NullExplorer)

type NullExplorer struct {
	SingleObjectiveAnnealableExplorer
}

func (nse *NullExplorer) Initialise() {
	nse.objectiveValue = float64(0)
}

func (nse *NullExplorer) WithName(name string) *NullExplorer {
	nse.SingleObjectiveAnnealableExplorer.WithName(name)
	return nse
}

func (nse *NullExplorer) SetObjectiveValue(temperature float64) {}
func (nse *NullExplorer) TryRandomChange(temperature float64)   {}
func (nse *NullExplorer) AcceptLastChange()                     {}
func (nse *NullExplorer) RevertLastChange()                     {}
