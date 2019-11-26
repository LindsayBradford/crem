// Copyright (c) 2019 Australian Rivers Institute.

package command

import "github.com/LindsayBradford/crem/pkg/attributes"

type Command interface {
	Do()
	Undo()
}

type BaseCommand struct {
	attributes.ContainedAttributes
	target interface{}
}

func (bc *BaseCommand) WithTarget(target interface{}) *BaseCommand {
	bc.target = target
	return bc
}

func (bc *BaseCommand) Target() interface{} {
	return bc.target
}

func (bc *BaseCommand) Do() {
	// deliberately does nothing
}

func (bc *BaseCommand) Undo() {
	// deliberately does nothing
}
