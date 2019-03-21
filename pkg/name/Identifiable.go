// Copyright (c) 2019 Australian Rivers Institute.

package name

// Identifiable offers an interface for anything needing an identifier
type Identifiable interface {
	SetId(title string)
	Id() string
}

// IdentifiableContainer is a struct offering a default implementation of Identifiable
type IdentifiableContainer struct {
	id string
}

func (n *IdentifiableContainer) Id() string {
	return n.id
}

func (n *IdentifiableContainer) SetId(id string) {
	n.id = id
}
