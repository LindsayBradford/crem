// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

var NULL_SOLUTION_EXPLORER = new(NullSolutionExplorer)

type NullSolutionExplorer struct {
	BaseSolutionExplorer
}

func (nse *NullSolutionExplorer) Initialise() {
	nse.objectiveValue = float64(0)
}

func (nse *NullSolutionExplorer) WithName(name string) *NullSolutionExplorer {
	nse.BaseSolutionExplorer.WithName(name)
	return nse
}

func (nse *NullSolutionExplorer) SetObjectiveValue(temperature float64) {}
func (nse *NullSolutionExplorer) TryRandomChange(temperature float64)   {}
func (nse *NullSolutionExplorer) AcceptLastChange()                     {}
func (nse *NullSolutionExplorer) RevertLastChange()                     {}
