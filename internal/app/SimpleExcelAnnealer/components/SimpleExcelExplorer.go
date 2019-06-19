// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/kirkpatrick"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
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

func (see *SimpleExcelExplorer) TryRandomChange() {
	see.Model().TryRandomChange()
	see.Explorer.AcceptOrRevertChange(see.AcceptLastChange, see.RevertLastChange)
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
	finishedIterationAttributes := see.EventAttributes(observer.FinishedIteration)

	return ExplorerData{
		Temperature:           finishedIterationAttributes.Value(explorer.Temperature).(float64),
		ChangeIsDesirable:     finishedIterationAttributes.Value(explorer.ChangeIsDesirable).(bool),
		ChangeAccepted:        finishedIterationAttributes.Value(explorer.ChangeAccepted).(bool),
		AcceptanceProbability: finishedIterationAttributes.Value(explorer.AcceptanceProbability).(float64),
	}
}
