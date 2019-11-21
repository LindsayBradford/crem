// Copyright (c) 2019 Australian Rivers Institute.

package command

func NewCompositeCommand() *CompositeCommand {
	newCommand := new(CompositeCommand)
	newCommand.composedCommands = make([]Command, 0)
	return newCommand
}

type CompositeCommand struct {
	composedCommands []Command
}

func (c *CompositeCommand) ComposedOf(commands ...Command) *CompositeCommand {
	c.Add(commands...)
	return c
}

func (c *CompositeCommand) Add(commands ...Command) {
	for _, command := range commands {
		c.composedCommands = append(c.composedCommands, command)
	}
}

func (c *CompositeCommand) Do() {
	for _, command := range c.composedCommands {
		command.Do()
	}
}

func (c *CompositeCommand) Undo() {
	for _, command := range c.composedCommands {
		command.Undo()
	}
}
