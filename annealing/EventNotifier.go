// Copyright (c) 2018 Australian Rivers Institute.

package annealing

import (
	"errors"
)

type EventNotifier interface {
	AddObserver(observer Observer) error
	Observers() []Observer
	NotifyObserversOfAnnealingEvent(annealer Annealer, eventType EventType)
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

func (notifier *SynchronousAnnealingEventNotifier) NotifyObserversOfAnnealingEvent(annealer Annealer, eventType EventType) {
	event := Event{EventType: eventType, Annealer: annealer}
	for _, currObserver := range notifier.observers {
		currObserver.ObserveAnnealingEvent(event)
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
	newEventChannel := make(chan Event)
	notifier.observerChannels = append(notifier.observerChannels, newEventChannel)

	go func() {
		for {
			newEvent := <-newEventChannel
			newObserver.ObserveAnnealingEvent(newEvent)
		}
	}()

	return nil
}

func (notifier *ConcurrentAnnealingEventNotifier) NotifyObserversOfAnnealingEvent(annealer Annealer, eventType EventType) {
	event := Event{EventType: eventType, Annealer: annealer}
	for _, observerChannel := range notifier.observerChannels {
		observerChannel <- event
	}
}
