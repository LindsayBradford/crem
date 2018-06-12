package Annealer

type AnnealerBuilder struct{}

func (builder *AnnealerBuilder) Build() (Annealer, error) {
	annealer := new(DefaultAnnealer)
	return annealer, nil
}
