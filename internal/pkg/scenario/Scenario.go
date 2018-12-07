// Copyright (c) 2018 Australian Rivers Institute.

package scenario

// Identifiable is an interface for anything needing a name
type Identifiable interface {
	ScenarioId() string
	SetScenarioId(id string)
}

// ContainedScenarioId is a struct offering a default implementation of Identifiable
type ContainedScenarioId struct {
	scenarioId string
}

func (n *ContainedScenarioId) ScenarioId() string {
	return n.scenarioId
}

func (n *ContainedScenarioId) SetScenarioId(id string) {
	n.scenarioId = id
}
