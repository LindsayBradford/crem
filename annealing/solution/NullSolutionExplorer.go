// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

var NULL_SOLUTION_EXPLORER = new(NullSolutionExplorer)

type NullSolutionExplorer struct {
	BaseSolutionExplorer
}

func (this *NullSolutionExplorer) Initialise() {
	this.objectiveValue = float64(0)
}

func (this *NullSolutionExplorer) SetObjectiveValue(temperature float64) {}
func (this *NullSolutionExplorer) TryRandomChange(temperature float64)   {}
func (this *NullSolutionExplorer) AcceptLastChange()                     {}
func (this *NullSolutionExplorer) RevertLastChange()                     {}
