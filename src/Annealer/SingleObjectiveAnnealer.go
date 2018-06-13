// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package Annealer

import "fmt"

type singleObjectiveAnnealer struct {
	abstractAnnealer
}

func (annealer *singleObjectiveAnnealer) Anneal() {
	fmt.Println("I'm a single-objectve annealer, annealing")
	fmt.Printf("Current Temperature: %f\n", annealer.Temperature())
	fmt.Printf("Iterations left: %d\n", annealer.IterationsLeft())
}
