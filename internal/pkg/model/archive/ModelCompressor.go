// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/pkg/archive"
	"github.com/LindsayBradford/crem/pkg/dominance"
)

type ModelCompressor struct{}

func (mc *ModelCompressor) Compress(model model.Model) *CompressedModelState {
	return &CompressedModelState{
		Variables: compressVariables(model),
		Actions:   compressActions(model),
	}
}

func compressVariables(model model.Model) dominance.Float64Vector {
	variableKeys := model.DecisionVariables().SortedKeys()
	compressedVariables := *dominance.NewFloat64(len(variableKeys))
	for index := range variableKeys {
		valueToCompress := model.DecisionVariable(variableKeys[index]).Value()
		compressedVariables[index] = valueToCompress
	}
	return compressedVariables
}

func compressActions(model model.Model) archive.BooleanArchive {
	actions := model.ManagementActions()
	compressedActions := *archive.New(len(actions))

	for index, action := range actions {
		compressedActions.SetValue(index, action.IsActive())
	}

	return compressedActions
}

func (mc *ModelCompressor) Decompress(condensedModelState *CompressedModelState, model model.Model) {
	for index := 0; index < condensedModelState.Actions.Len(); index++ {
		compressedActionState := condensedModelState.Actions.Value(index)
		model.SetManagementAction(index, compressedActionState)
	}
}
