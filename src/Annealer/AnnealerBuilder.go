// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package Annealer

type AnnealerBuilder struct{}

func (builder *AnnealerBuilder) Build() (Annealer, error) {
	annealer := new(defaultAnnealer)
	return annealer, nil
}
