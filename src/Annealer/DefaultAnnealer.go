package Annealer

import "fmt"

type Annealer interface {
	Anneal()
}

type DefaultAnnealer struct {
}

func (annealer *DefaultAnnealer) Anneal() {
	fmt.Println("I'm an annealer, annealing")
}
