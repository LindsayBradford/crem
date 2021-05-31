// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/pkg/archive"
	"github.com/LindsayBradford/crem/pkg/dominance"
	"github.com/LindsayBradford/crem/pkg/name"
)

type CompressedModelState struct {
	name.IdentifiableContainer

	Variables dominance.Float64Vector
	Actions   archive.BooleanArchive
}

func (c *CompressedModelState) MatchesStateOf(model model.Model) bool {
	return c.variablesMatch(model) && c.actionsMatch(model)
}

func (c *CompressedModelState) variablesMatch(model model.Model) bool {
	variableKeys := model.DecisionVariables().SortedKeys()
	for index := range variableKeys {
		if !c.variableValuesMatch(index, model, variableKeys) {
			return false
		}
	}
	return true
}

func (c *CompressedModelState) VariableDifferences(otherState *CompressedModelState) []float64 {
	differences := make([]float64, len(c.Variables))
	for index := range c.Variables {
		differences[index] = c.Variables[index] - otherState.Variables[index]
	}
	return differences
}

func (c *CompressedModelState) variableValuesMatch(index int, model model.Model, variableKeys []string) bool {
	return c.Variables[index] == model.DecisionVariable(variableKeys[index]).Value()
}

func (c *CompressedModelState) actionsMatch(model model.Model) bool {
	for index := range model.ManagementActions() {
		if !c.actionValuesMatch(index, model) {
			return false
		}
	}
	return true
}

func (c *CompressedModelState) actionValuesMatch(index int, model model.Model) bool {
	return c.Actions.Value(index) == model.ManagementActions()[index].IsActive()
}

func (c *CompressedModelState) IsEquivalentTo(otherSate *CompressedModelState) bool {
	return c.Actions.IsEquivalentTo(&otherSate.Actions)
}

func (c *CompressedModelState) Encoding() string {
	return c.Actions.Encoding()
}

func (c *CompressedModelState) Decode(encoding string) error {
	return c.Actions.Decode(encoding)
}
