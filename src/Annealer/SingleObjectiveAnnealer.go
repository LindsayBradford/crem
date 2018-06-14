// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package Annealer

import "fmt"

type singleObjectiveAnnealer struct {
	annealerBase
}

func (this *singleObjectiveAnnealer) Anneal() {
	fmt.Println("I'm a single-objectve annealer")
	this.annealerBase.Anneal()
}
