package Annealer

import "fmt"

type Annealer struct{}

func (annealer *Annealer) Anneal() {
	fmt.Println("I'm an annealer, annealing")
}
