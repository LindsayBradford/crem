// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging


import (
	. "github.com/LindsayBradford/crm/annealing"
)

type AnnealingLogger struct {
	logHandler LogHandler
}

func (this *AnnealingLogger) ObserveAnnealingEvent(event AnnealingEvent) {
	//Deiberately does nothing but ensure it matches the AnnealingObserver interface
}