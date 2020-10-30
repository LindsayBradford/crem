// Copyright (c) 2019 Australian Rivers Institute.

package observer

import (
	"errors"
)

type EventNotifier interface {
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

type ConcurrentAnnealingEventNotifier struct {
	observers        []Observer
	observerChannels []chan Event
}

func (notifier *ConcurrentAnnealingEventNotifier) Observers() []Observer {
	if len(notifier.observers) == 0 {
		return nil
	}
	return notifier.observers
}

func (notifier *ConcurrentAnnealingEventNotifier) AddObserver(newObserver Observer) error {
	if newObserver == nil {
		return errors.New("invalid attempt to add non-existent observer to annealing event notifier")
	}

	notifier.observers = append(notifier.observers, newObserver)
	return notifier.addObserverToChannels(newObserver)
}

func (notifier *ConcurrentAnnealingEventNotifier) AddObserverAsFirst(newObserver Observer) error {
	if newObserver == nil {
		return errors.New("invalid attempt to add non-existent observer to annealing event notifier")
	}

	notifier.observers = append([]Observer{newObserver}, notifier.observers...)
	return notifier.addObserverToChannels(newObserver)
}

func (notifier *ConcurrentAnnealingEventNotifier) addObserverToChannels(newObserver Observer) error {
	newEventChannel := make(chan Event)
	notifier.observerChannels = append(notifier.observerChannels, newEventChannel)

	go func() {
		for {
			newEvent := <-newEventChannel
			newObserver.ObserveEvent(newEvent)
		}
	}()

	return nil
}

func (notifier *ConcurrentAnnealingEventNotifier) NotifyObserversOfEvent(event Event) {
	for _, observerChannel := range notifier.observerChannels {
		observerChannel <- event
	}
}
