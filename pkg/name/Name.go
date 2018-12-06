// Copyright (c) 2018 Australian Rivers Institute.

package name

// Nameable is an interface for anything needing a name
type Nameable interface {
	Name() string
	SetName(name string)
}

// Named is a struct offering a default implementation of Nameable
type Named struct {
	name string
}

func (n *Named) Name() string {
	return n.name
}

func (n *Named) SetName(name string) {
	n.name = name
}
