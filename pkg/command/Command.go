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

func (c *BaseCommand) WithTarget(target interface{}) *BaseCommand {
	c.target = target
	return c
}

func (c *BaseCommand) Target() interface{} {
	return c.target
}

func (c *BaseCommand) Reset() {
	c.status = UnDone
}

func (c *BaseCommand) Do() CommandStatus {
	if c.status == UnDone {
		c.status = Done
		return c.status
	}
	return NoChange
}

func (c *BaseCommand) Undo() CommandStatus {
	if c.status == Done {
		c.status = UnDone
		return c.status
	}
	return NoChange
}
