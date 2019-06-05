// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/pkg/archive"
	"github.com/LindsayBradford/crem/pkg/dominance"
)

type Archivist struct{}

func (a *Archivist) Store(model model.Model) *ModelArchive {
	vector := a.buildVariableVector(model)
	actions := a.buildActionArchive(model)

	return new(ModelArchive).WithVariables(vector).WithActions(actions)
}

func (a *Archivist) buildVariableVector(model model.Model) dominance.Float64Vector {
	variableKeys := model.DecisionVariables().SortedKeys()
	vector := *dominance.NewFloat64(len(variableKeys))
	for index := range variableKeys {
		vector[index] = model.DecisionVariable(variableKeys[index]).Value()
	}
	return vector
}

func (a *Archivist) buildActionArchive(model model.Model) archive.BooleanArchive {
	actions := model.ManagementActions()
	archive := *archive.New(len(actions))

	for index, action := range actions {
		archive.SetValue(index, action.IsActive())
	}

	return archive
}

func (a *Archivist) Retrieve(archive *ModelArchive, model model.Model) {
	actionArchive := archive.Actions()
	for index := 0; index < actionArchive.Len(); index++ {
		model.SetManagementActionUnobserved(index, actionArchive.Value(index))
	}
}
