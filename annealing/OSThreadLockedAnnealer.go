// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealing

import (
	"runtime"
)

type OSThreadLockedAnnealer struct {
	ElapsedTimeTrackingAnnealer
}

func (this *OSThreadLockedAnnealer) Anneal() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	this.ElapsedTimeTrackingAnnealer.Anneal()
}
