// Copyright (c) 2018 Australian Rivers Institute.

// name package offers a suite of interfaces and structs, allowing uniform naming/identification of entities.
package name

// Nameable is an interface for anything needing a name
type Nameable interface {
	Name() string
	SetName(name string)
}

// ContainedLogger offers a default implementation of Nameable
type NameContainer struct {
	name string
}

func (n *NameContainer) Name() string {
	return n.name
}

func (n *NameContainer) SetName(name string) {
	n.name = name
}
