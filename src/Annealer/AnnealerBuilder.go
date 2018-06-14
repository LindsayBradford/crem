// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package Annealer

type AnnealerBuilder struct {
	annealer Annealer
}

func (builder *AnnealerBuilder) SingleObjectiveAnnealer() *AnnealerBuilder {
	builder.annealer = &singleObjectiveAnnealer{}
	builder.annealer.Initialise()
	return builder
}

func (builder *AnnealerBuilder) WithStartingTemperature(temperature float64) *AnnealerBuilder {
	annealer := builder.annealer
	annealer.setTemperature(temperature)
	return builder
}

func (builder *AnnealerBuilder) WithCoolingFactor(coolingFactor float64) *AnnealerBuilder {
	annealer := builder.annealer
	annealer.setCoolingFactor(coolingFactor)
	return builder
}

func (builder *AnnealerBuilder) WithMaxIterations(iterations uint) *AnnealerBuilder {
	annealer := builder.annealer
	annealer.setMaxIterations(iterations)
	return builder
}

func (builder *AnnealerBuilder) Build() Annealer {
	return builder.annealer
}
