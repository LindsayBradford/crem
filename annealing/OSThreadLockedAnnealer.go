// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealing

import (
	"runtime"
)

type OSThreadLockedAnnealer struct {
	ElapsedTimeTrackingAnnealer
}

func (annealer *OSThreadLockedAnnealer) Initialise() {
	annealer.ElapsedTimeTrackingAnnealer.Initialise()
	annealer.SimpleAnnealer.SetTitle("Operating System Thread-Locked Annealer")
}

func (annealer *OSThreadLockedAnnealer) Anneal() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	annealer.ElapsedTimeTrackingAnnealer.Anneal()
}
