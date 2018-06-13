package Annealer

import (
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	"testing"
)

func TestBuild(t *testing.T) {

	const expectedTemperature float64 = 1000.0
	const expectedIterations uint = 2000

	builder := new(AnnealerBuilder)

	annealer, _ := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(expectedTemperature).
		WithIterations(expectedIterations).
		Build()

	AssertThat(t, annealer.Temperature(), is.EqualTo(expectedTemperature).Reason("Build() Temperature"))
	AssertThat(t, annealer.IterationsLeft(), is.EqualTo(expectedIterations).Reason("Build() IterationsLeft"))
}
