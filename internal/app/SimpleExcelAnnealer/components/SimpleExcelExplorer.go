// Copyright (c) 2018 Australian Rivers Institute.

package components

import "github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/kirkpatrick"

func NewSimpleExcelExplorer() *SimpleExcelExplorer {
	explorer := new(SimpleExcelExplorer)
	explorer.Explorer = *kirkpatrick.New()
	return explorer
}

type SimpleExcelExplorer struct {
	kirkpatrick.Explorer
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
