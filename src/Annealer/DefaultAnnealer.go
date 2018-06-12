package Annealer

import "fmt"

type DefaultAnnealer struct {
}

func (annealer *DefaultAnnealer) Anneal() {
	fmt.Println("I'm an annealer, annealing")
}
