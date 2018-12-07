// Copyright (c) 2018 Australian Rivers Institute.

package name

// Nameable is an interface for anything needing a name
type Nameable interface {
	Name() string
	SetName(name string)
}

// ContainedName is a struct offering a default implementation of Nameable
type ContainedName struct {
	name string
}

func (n *ContainedName) Name() string {
	return n.name
}

func (n *ContainedName) SetName(name string) {
	n.name = name
}
