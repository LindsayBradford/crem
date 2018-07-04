// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package objectives

var NULL_OBJECTIVE_MANAGER = new(NullObjectiveManager)

type NullObjectiveManager struct {
	BaseObjectiveManager
}

func (this *NullObjectiveManager) Initialise() {
	this.objectiveValue = float64(0)
}

func (this *NullObjectiveManager) TryRandomChange(temperature float64) {}
func (this *NullObjectiveManager) AcceptLastChange()                   {}
func (this *NullObjectiveManager) RevertLastChange()                   {}

