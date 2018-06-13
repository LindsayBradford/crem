// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package Annealer

type AnnealerBuilder struct {
	annealer Annealer
}

func (builder *AnnealerBuilder) SingleObjectiveAnnealer() *AnnealerBuilder {
	builder.annealer = &singleObjectiveAnnealer{}
	return builder
}

func (builder *AnnealerBuilder) WithStartingTemperature(temperature float64) *AnnealerBuilder {
	annealer := builder.annealer
	annealer.setTemperature(temperature)
	return builder
}

func (builder *AnnealerBuilder) WithIterations(iterations uint) *AnnealerBuilder {
	annealer := builder.annealer
	annealer.setIterationsLeft(iterations)
	return builder
}

func (builder *AnnealerBuilder) Build() (Annealer, error) {
	finalisedAnnealer := builder.annealer
	return finalisedAnnealer, nil
}
