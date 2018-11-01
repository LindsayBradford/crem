// Copyright (c) 2018 Australian Rivers Institute.

package filters

import "github.com/LindsayBradford/crem/annealing"

// PercentileOfIterationsPerAnnealingFilter filters FinishedIteration Annealing Event instances at a rate of 1 every
// percentile number of iterations received. . StartedIteration events are completely filtered out. All other event types are allowed through to the LogHandler.
type PercentileOfIterationsPerAnnealingFilter struct {
	maxIterations      uint64
	percentileToReport float64
	iterationModulo    uint64
	generatesEvents    bool
}

// WithPercentileOfIterations defines the number of Annealing Iteration Event instances to report over the entire run.
func (m *PercentileOfIterationsPerAnnealingFilter) WithPercentileOfIterations(percentile float64) *PercentileOfIterationsPerAnnealingFilter {
	m.percentileToReport = percentile / 100
	m.deriveModuloIfPossible()
	return m
}

// WithPercentileOfIterations defines the number of Annealing Iteration Event instances to report over the entire run.
func (m *PercentileOfIterationsPerAnnealingFilter) WithMaxIterations(maxIterations uint64) *PercentileOfIterationsPerAnnealingFilter {
	m.maxIterations = maxIterations
	m.deriveModuloIfPossible()
	return m
}

func (m *PercentileOfIterationsPerAnnealingFilter) deriveModuloIfPossible() {
	m.generatesEvents = false
	if m.maxIterations > 0 && m.percentileToReport > 0 {
		m.generatesEvents = true
		if m.percentileToReport > 1 {
			m.percentileToReport = 1
		}
		if m.percentileToReport == 1 {
			m.iterationModulo = 1
		} else {
			m.iterationModulo = uint64((float64)(m.maxIterations) * m.percentileToReport)
		}
	}
}

// ShouldFilter filters only FinishedIteration Event instances, and fully filters out all StartedIteration
// events. Every FinishedIteration events received on the specified percentile boundary, one event is allowed through
// to the LogHandler.
func (m *PercentileOfIterationsPerAnnealingFilter) ShouldFilter(event annealing.Event) bool {
	if event.EventType != annealing.StartedIteration && event.EventType != annealing.FinishedIteration {
		return false
	}

	annealer := event.Annealer

	if m.generatesEvents && event.EventType == annealing.FinishedIteration &&
		annealer.CurrentIteration()%m.iterationModulo == 0 {
		return false
	}

	return true
}
