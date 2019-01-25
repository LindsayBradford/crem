// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/kirkpatrick"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
)

func NewSimpleExcelExplorer() *SimpleExcelExplorer {
	explorer := new(SimpleExcelExplorer)
	explorer.Explorer = *kirkpatrick.New()
	return explorer
}

type SimpleExcelExplorer struct {
	kirkpatrick.Explorer
}

func (see *SimpleExcelExplorer) WithName(name string) *SimpleExcelExplorer {
	see.Explorer.WithName(name)
	return see
}

func (see *SimpleExcelExplorer) WithModel(model model.Model) *SimpleExcelExplorer {
	see.Explorer.WithModel(model)
	return see
}

func (see *SimpleExcelExplorer) WithScenarioId(id string) *SimpleExcelExplorer {
	see.Explorer.WithScenarioId(id)
	return see
}

func (see *SimpleExcelExplorer) WithParameters(params parameters.Map) *SimpleExcelExplorer {
	see.Explorer.WithParameters(params)
	return see
}

func (see *SimpleExcelExplorer) DeepClone() explorer.Explorer {
	clone := *see
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	modelClone := see.Model().DeepClone()
	clone.SetModel(modelClone)
	return &clone
}

func (see *SimpleExcelExplorer) TryRandomChange(temperature float64) {
	see.Model().TryRandomChange()
	see.Explorer.AcceptOrRevertChange(temperature, see.AcceptLastChange, see.RevertLastChange)
}

func (see *SimpleExcelExplorer) AcceptLastChange() {
	see.Explorer.AcceptLastChange()
	see.simpleExcelModel().SetExplorerData(see.buildExplorerData())
}

func (see *SimpleExcelExplorer) RevertLastChange() {
	see.Explorer.RevertLastChange()
	see.simpleExcelModel().SetExplorerData(see.buildExplorerData())
}

func (see *SimpleExcelExplorer) simpleExcelModel() *SimpleExcelModel {
	if model, isValid := see.Model().(*SimpleExcelModel); isValid {
		return model
	}
	panic("cant cast model to a SimpleExcelModel")
}

func (see *SimpleExcelExplorer) buildExplorerData() ExplorerData {
	data := ExplorerData{
		Temperature:           see.Temperature(),
		ChangeIsDesirable:     see.ChangeIsDesirable(),
		ChangeAccepted:        see.ChangeAccepted(),
		AcceptanceProbability: see.AcceptanceProbability(),
	}
	return data
}
