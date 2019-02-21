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
	EventType  EventType
	attributes attributes.Attributes
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

func (e *Event) Note() string {
	return e.Attribute(noteTextKey).(string)
}

func (e *Event) WithAttribute(name string, value interface{}) *Event {
	e.AddAttribute(name, value)
	return e
}

func (e *Event) Attribute(name string) interface{} {
	return e.attributes.Value(name)
}

func (e *Event) AddAttribute(name string, value interface{}) {
	newEntry := attributes.NameValuePair{Name: name, Value: value}
	e.attributes = append(e.attributes, newEntry)
}

type EventType int

const (
	InvalidEvent EventType = iota
	StartedAnnealing
	StartedIteration
	FinishedIteration
	FinishedAnnealing
	ManagementAction
	DecisionVariable
	Note
)

func (eventType EventType) String() string {
	labels := [...]string{
		"InvalidEvent",
		"StartedAnnealing",
		"StartedIteration",
		"FinishedIteration",
		"FinishedAnnealing",
		"ManagementAction",
		"DecisionVariable",
		"Note",
	}

	if eventType < StartedAnnealing || eventType > Note {
		return labels[InvalidEvent]
	}

	return labels[eventType]
}
