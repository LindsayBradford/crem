// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package Annealer

import "fmt"

type defaultAnnealer struct {
}

func (annealer *defaultAnnealer) Anneal() {
	fmt.Println("I'm an annealer, annealing")
}
