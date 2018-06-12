// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package Annealer

import "fmt"

type defaultAnnealer struct {
	abstractAnnealer
}

func (annealer *defaultAnnealer) Anneal() {
	fmt.Println("I'm an annealer, annealing")
	fmt.Printf("Current Temperature: %f\n", annealer.Temperature())
	fmt.Printf("Iterations left: %d\n", annealer.IterationsLeft())
}
