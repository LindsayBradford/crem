// Copyright (c) 2018 Australian Rivers Institute.

package api

import "github.com/LindsayBradford/crem/server/job"

type JobQueue struct {
	job.Queue
}

func (jq *JobQueue) Initialise() *JobQueue {
	jq.Queue.Initialise()
	return jq
}

func (jq *JobQueue) WithQueueLength(length uint64) *JobQueue {
	jq.Queue.WithQueueLength(length)
	return jq
}
