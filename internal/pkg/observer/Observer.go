// Copyright (c) 2019 Australian Rivers Institute.

package observer

import "github.com/LindsayBradford/crem/pkg/attributes"

type Observer interface {
	ObserveEvent(event Event)
}

const (
	_           = iota
	idKey       = "Id"
	sourceKey   = "Source"
	noteTextKey = "Note"
)

type Event struct {
	EventType EventType
	attributes.ContainedAttributes
}

func NewEvent(eventType EventType) *Event {
	newEvent := new(Event)
	newEvent.EventType = eventType
	return newEvent
}

func (e *Event) WithId(id string) *Event {
	e.AddAttribute(idKey, id)
	return e
}

func (e *Event) Id() string {
	return e.Attribute(idKey).(string)
}

func (e *Event) WithSource(source interface{}) *Event {
	e.AddAttribute(sourceKey, source)
	return e
}

func (e *Event) Source() interface{} {
	return e.Attribute(sourceKey)
}

func (e *Event) WithNote(text string) *Event {
	e.AddAttribute(noteTextKey, text)
	return e
}

func (e *Event) ReplaceNote(text string) *Event {
	e.ReplaceAttribute(noteTextKey, text)
	return e
}

func (e *Event) Note() string {
	note := e.Attribute(noteTextKey)
	if note == nil {
		return ""
	}
	return note.(string)
}

func (e *Event) HasNote() bool {
	return e.Attribute(noteTextKey) != nil
}

func (e *Event) WithAttribute(name string, value interface{}) *Event {
	e.AddAttribute(name, value)
	return e
}

func (e *Event) ReplacingAttribute(name string, value interface{}) *Event {
	e.ContainedAttributes.ReplaceAttribute(name, value)
	return e
}

func (e *Event) JoiningAttributes(newAttributes attributes.Attributes) *Event {
	e.ContainedAttributes.JoiningAttributes(newAttributes)
	return e
}

type EventType int

const (
	InvalidEvent EventType = iota
	StartedAnnealing
	StartedIteration
	Explorer
	Model
	ManagementAction
	DecisionVariable
	FinishedIteration
	FinishedAnnealing
	Note
)

func (eventType EventType) String() string {
	labels := [...]string{
		"InvalidEvent",
		"StartedAnnealing",
		"StartedIteration",
		"Explorer",
		"Model",
		"ManagementAction",
		"DecisionVariable",
		"FinishedIteration",
		"FinishedAnnealing",
		"Note",
	}

	if eventType < StartedAnnealing || eventType > Note {
		return labels[InvalidEvent]
	}

	return labels[eventType]
}

func (eventType EventType) IsAnnealingState() bool {
	return eventType >= StartedAnnealing && eventType <= FinishedAnnealing
}

func (eventType EventType) IsAnnealingIterationState() bool {
	return eventType >= StartedIteration && eventType <= FinishedIteration
}
func (eventType EventType) IsModelState() bool {
	return eventType >= Model && eventType <= DecisionVariable
}
