// Copyright (c) 2019 Australian Rivers Institute.

package command

import (
	"github.com/LindsayBradford/crem/pkg/attributes"
)

type CommandStatus int

const UnDone CommandStatus = 0
const Done CommandStatus = 1
const NoChange CommandStatus = 2

type Command interface {
	Do() CommandStatus
	Undo() CommandStatus
	Reset()
}

type BaseCommand struct {
	attributes.ContainedAttributes
	target interface{}
	status CommandStatus
}

func (bc *BaseCommand) WithTarget(target interface{}) *BaseCommand {
	bc.target = target
	return bc
}

func (bc *BaseCommand) Target() interface{} {
	return bc.target
}

func (bc *BaseCommand) Reset() {
	bc.status = UnDone
}

func (bc *BaseCommand) Do() CommandStatus {
	if bc.status == UnDone {
		bc.status = Done
		return bc.status
	}
	return NoChange
}

func (bc *BaseCommand) Undo() CommandStatus {
	if bc.status == Done {
		bc.status = UnDone
		return bc.status
	}
	return NoChange
}
