// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import "fmt"

type singleObjectiveAnnealer struct {
	annealerBase
}

func (this *singleObjectiveAnnealer) Anneal() {
	this.notifyObserversWith("I'm a single-objective annealer")
	this.annealerBase.Anneal()
	this.notifyObserversWith(this.generateElapsedTimeString())
}

func (this *singleObjectiveAnnealer) generateElapsedTimeString() string {
	return fmt.Sprintf("Total elapsed time of annealing = [%v]", this.ElapsedTime())
}
