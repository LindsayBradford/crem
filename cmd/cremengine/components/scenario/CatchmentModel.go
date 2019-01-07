// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/dumb"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
)

type CatchmentModel struct {
	dumb.Model
	parameters CatchmentParameters
}

func NewCatchmentModel() *CatchmentModel {
	newModel := new(CatchmentModel)
	newModel.SetName("CatchmentModel")
	newModel.Model = *dumb.New()

	newModel.parameters.Initialise()

	return newModel
}

func (cm *CatchmentModel) WithName(name string) *CatchmentModel {
	cm.SetName(name)
	return cm
}

func (cm *CatchmentModel) WithParameters(params parameters.Map) *CatchmentModel {
	cm.SetParameters(params)
	return cm
}

func (cm *CatchmentModel) SetParameters(params parameters.Map) error {
	cm.parameters.Merge(params)

	return cm.parameters.ValidationErrors()
}

func (cm *CatchmentModel) ParameterErrors() error {
	return cm.parameters.ValidationErrors()
}

func (cm *CatchmentModel) Initialise() {
	cm.Model.Initialise()

	initialValue := cm.parameters.GetFloat64(dumb.InitialObjectiveValue)
	cm.Model.SetDecisionVariable("ObjectiveValue", initialValue)
}

func (cm *CatchmentModel) TearDown() {
	cm.Model.TearDown()
}

func (cm *CatchmentModel) TryRandomChange() {
	// TODO: randomly choose a management action to toggle.
	cm.Model.TryRandomChange()
}

func (cm *CatchmentModel) DeepClone() model.Model {
	clone := *cm
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
