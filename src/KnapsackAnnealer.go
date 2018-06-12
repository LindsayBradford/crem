// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package main

import "Annealer"

func main() {
	builder := new(Annealer.AnnealerBuilder)
	annealer, _ := builder.Build()
	annealer.Anneal()
}
