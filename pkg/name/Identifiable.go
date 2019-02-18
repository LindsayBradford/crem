// Copyright (c) 2019 Australian Rivers Institute.

package name

type Identifiable interface {
	SetId(title string)
	Id() string
}

// ContainedIdentifier is a struct offering a default implementation of Identifiable
type ContainedIdentifier struct {
	id string
}

func (n *ContainedIdentifier) Id() string {
	return n.id
}

func (n *ContainedIdentifier) SetId(id string) {
	n.id = id
}
