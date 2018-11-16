// Copyright (c) 2018 Australian Rivers Institute.

package job

const defaultQueueLength = 1
const unspecifiedQueueLength = 0

type Queue struct {
	Jobs        chan *Job      `json:"-"`
	JobFunction func(job *Job) `json:"-"`
}

func (jq *Queue) Initialise() *Queue {
	jq.Jobs = make(chan *Job, defaultQueueLength)
	return jq
}

func (jq *Queue) WithQueueLength(length uint64) *Queue {
	if length > unspecifiedQueueLength {
		jq.Jobs = make(chan *Job, length)
	}
	return jq
}

type EnqueuedStatus bool

const EnqueueSucceeded EnqueuedStatus = true
const EnqueueFailed EnqueuedStatus = false

func (jq *Queue) Enqueue(newJob *Job) EnqueuedStatus {
	select {
	case jq.Jobs <- newJob:
		return EnqueueSucceeded
	default:
		return EnqueueFailed
	}
}

func (jq *Queue) Start() {
	for {
		job := <-jq.Jobs
		jq.JobFunction(job)
	}
}
