// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

var NULL_SOLUTION_TOURER = new(NullSolutionTourer)

type NullSolutionTourer struct {
	BaseSolutionTourer
}

func (this *NullSolutionTourer) Initialise() {
	this.objectiveValue = float64(0)
}

func (this *NullSolutionTourer) SetObjectiveValue(temperature float64) {}
func (this *NullSolutionTourer) TryRandomChange(temperature float64)   {}
func (this *NullSolutionTourer) AcceptLastChange()                     {}
func (this *NullSolutionTourer) RevertLastChange()                     {}

