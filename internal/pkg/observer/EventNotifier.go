// Copyright (c) 2019 Australian Rivers Institute.

package observer

import (
	"errors"
)

type EventNotifier interface {
	HasObservers() bool
	AddObserver(observer Observer) error
	AddObserverAsFirst(observer Observer) error
	Observers() []Observer
	NotifyObserversOfEvent(event Event)
}

// EventNotifierContainer  defines an interface embedding an EventNotifier
type EventNotifierContainer interface {
	EventNotifier() EventNotifier
	SetEventNotifier(notifier EventNotifier) error
}

// ContainedEventNotifier offers a struct implementing the EventNotifierContainer interface.
type ContainedEventNotifier struct {
	notifier EventNotifier
}

func (cen *ContainedEventNotifier) EventNotifier() EventNotifier {
	return cen.notifier
}

func (cen *ContainedEventNotifier) SetEventNotifier(notifier EventNotifier) error {
	if notifier == nil {
		return errors.New("invalid attempt to set event notifier to nil value")
	}
	cen.notifier = notifier
	return nil
}

type SynchronousAnnealingEventNotifier struct {
	observers []Observer
}

func (notifier *SynchronousAnnealingEventNotifier) HasObservers() bool {
	return len(notifier.observers) > 0
}

func (notifier *SynchronousAnnealingEventNotifier) Observers() []Observer {
	if len(notifier.observers) == 0 {
		return nil
	}
	return notifier.observers
}

func (notifier *SynchronousAnnealingEventNotifier) AddObserver(newObserver Observer) error {
	if newObserver == nil {
		return errors.New("invalid attempt to add non-existent observer to annealing event notifier")
	}
	notifier.observers = append(notifier.observers, newObserver)
	return nil
}

func (notifier *SynchronousAnnealingEventNotifier) AddObserverAsFirst(newObserver Observer) error {
	if newObserver == nil {
		return errors.New("invalid attempt to add non-existent observer to annealing event notifier")
	}
	notifier.observers = append([]Observer{newObserver}, notifier.observers...)
	return nil
}

func (notifier *SynchronousAnnealingEventNotifier) NotifyObserversOfEvent(event Event) {
	for _, currObserver := range notifier.observers {
		currObserver.ObserveEvent(event)
	}
}
