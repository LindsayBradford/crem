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
	observers        []AnnealingObserver
}

func (this *SynchronousAnnealingEventNotifier) Observers() []AnnealingObserver {
	if (len(this.observers) == 0) {
		return nil
	}
	return this.observers
}

func (this *SynchronousAnnealingEventNotifier) AddObserver(newObserver AnnealingObserver) error {
	if newObserver == nil {
		return errors.New("Invalid attempt to add non-existant observer to annealing event notifier")
	}
	this.observers = append(this.observers, newObserver)
	return nil
}

func (this *SynchronousAnnealingEventNotifier) NotifyObserversOfAnnealingEvent(annealer Annealer, eventType AnnealingEventType) {
	event := AnnealingEvent{EventType: eventType, Annealer: annealer}
	for _, currObserver := range this.observers {
		currObserver.ObserveAnnealingEvent(event)
	}
}

type ChanneledAnnealingEventNotifier struct {
	observers        []AnnealingObserver
	observerChannels []chan AnnealingEvent
}

func (this *ChanneledAnnealingEventNotifier) Observers() []AnnealingObserver {
	if (len(this.observers) == 0) {
		return nil
	}
	return this.observers
}

func (this *ChanneledAnnealingEventNotifier) AddObserver(newObserver AnnealingObserver) error {
	if newObserver == nil {
		return errors.New("Invalid attempt to add non-existant observer to annealing event notifier")
	}

	this.observers = append(this.observers, newObserver)
	newEventChannel := make(chan AnnealingEvent)
	this.observerChannels = append(this.observerChannels, newEventChannel)

	go func() {
		for {
			newEvent := <- newEventChannel
			newObserver.ObserveAnnealingEvent(newEvent)
		}
	}()

	return nil
}

func (this *ChanneledAnnealingEventNotifier) NotifyObserversOfAnnealingEvent(annealer Annealer, eventType AnnealingEventType) {
	event := AnnealingEvent{ EventType: eventType, Annealer:  annealer}
	for _, observerChannel := range this.observerChannels {
		observerChannel <- event
	}
}