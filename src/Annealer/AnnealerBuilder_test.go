package Annealer

import "testing"

func TestBuild(t *testing.T) {

	const expectedTemperature = 1000.0
	const expectedIterations = 5000

	builder := new(AnnealerBuilder)

	annealer, _ := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(expectedTemperature).
		WithIterations(expectedIterations).
		Build()

	if annealer.Temperature() != expectedTemperature {
		t.Errorf("Expecting temperature of %f, got %f", expectedTemperature, annealer.Temperature())
	}

	if annealer.IterationsLeft() != expectedIterations {
		t.Errorf("Expecting iterations left of %d, got %d", expectedIterations, annealer.IterationsLeft())
	}
}
