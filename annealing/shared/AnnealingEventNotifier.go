// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"errors"
)

type AnnealingEventNotifier interface {
	AddObserver(observer AnnealingObserver) error
	Observers() []AnnealingObserver
	NotifyObserversOfAnnealingEvent(annealer Annealer, eventType AnnealingEventType)
}

type SynchronousAnnealingEventNotifier struct {
	observers []AnnealingObserver
}

func (notifier *SynchronousAnnealingEventNotifier) Observers() []AnnealingObserver {
	if len(notifier.observers) == 0 {
		return nil
	}
	return notifier.observers
}

func (notifier *SynchronousAnnealingEventNotifier) AddObserver(newObserver AnnealingObserver) error {
	if newObserver == nil {
		return errors.New("invalid attempt to add non-existent observer to annealing event notifier")
	}
	notifier.observers = append(notifier.observers, newObserver)
	return nil
}

func (notifier *SynchronousAnnealingEventNotifier) NotifyObserversOfAnnealingEvent(annealer Annealer, eventType AnnealingEventType) {
	event := AnnealingEvent{EventType: eventType, Annealer: annealer}
	for _, currObserver := range notifier.observers {
		currObserver.ObserveAnnealingEvent(event)
	}
}

type ChanneledAnnealingEventNotifier struct {
	observers        []AnnealingObserver
	observerChannels []chan AnnealingEvent
}

func (notifier *ChanneledAnnealingEventNotifier) Observers() []AnnealingObserver {
	if len(notifier.observers) == 0 {
		return nil
	}
	return notifier.observers
}

func (notifier *ChanneledAnnealingEventNotifier) AddObserver(newObserver AnnealingObserver) error {
	if newObserver == nil {
		return errors.New("invalid attempt to add non-existent observer to annealing event notifier")
	}

	notifier.observers = append(notifier.observers, newObserver)
	newEventChannel := make(chan AnnealingEvent)
	notifier.observerChannels = append(notifier.observerChannels, newEventChannel)

	go func() {
		for {
			newEvent := <-newEventChannel
			newObserver.ObserveAnnealingEvent(newEvent)
		}
	}()

	return nil
}

func (notifier *ChanneledAnnealingEventNotifier) NotifyObserversOfAnnealingEvent(annealer Annealer, eventType AnnealingEventType) {
	event := AnnealingEvent{EventType: eventType, Annealer: annealer}
	for _, observerChannel := range notifier.observerChannels {
		observerChannel <- event
	}
}
