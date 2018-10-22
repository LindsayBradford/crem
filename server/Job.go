// Copyright (c) 2018 Australian Rivers Institute.

package server

import "github.com/nu7hatch/gouuid"

const creationTimeKey = "CreationTime"
const completionTimeKey = "CompletedTime"

const statusKey = "Status"

const defaultQueueLength = 1 // TODO: Make this configurable.

type JobQueue struct {
	Jobs        chan *Job      `json:"-"`
	JobFunction func(job *Job) `json:"-"`
}

func (jq *JobQueue) Initialise() *JobQueue {
	jq.Jobs = make(chan *Job, defaultQueueLength)
	return jq
}

type JobEnqueuedStatus bool

const JobEnqueueSucceeded JobEnqueuedStatus = true
const JobEnqueueFailed JobEnqueuedStatus = false

func (jq *JobQueue) Enqueue(newJob *Job) JobEnqueuedStatus {
	select {
	case jq.Jobs <- newJob:
		return JobEnqueueSucceeded
	default:
		return JobEnqueueFailed
	}
}

func (jq *JobQueue) Start() {
	for {
		job := <-jq.Jobs
		jq.JobFunction(job)
	}
}

type JobId string
type JobStatus JobAttributeKey
type JobAttributeKey string

const (
	JobCreated   JobStatus = "CREATED"
	JobCompleted JobStatus = "COMPLETED"
	JobInvalid   JobStatus = "INVALID"
)

type Job struct {
	Id               JobId
	Attributes       map[JobAttributeKey]interface{}
	HiddenAttributes map[JobAttributeKey]interface{} `json:"-"`
}

func (j *Job) Initialise() *Job {
	j.createNewJobID()
	j.makeAttributeMaps()
	j.recordCreationAttributes()
	return j
}

func (j *Job) createNewJobID() {
	newId, _ := uuid.NewV4()
	j.Id = JobId(newId.String())
}

func (j *Job) makeAttributeMaps() {
	j.Attributes = make(map[JobAttributeKey]interface{}, 0)
	j.HiddenAttributes = make(map[JobAttributeKey]interface{}, 0)
}

func (j *Job) recordCreationAttributes() {
	j.recordCreationTime()
	j.SetStatus(JobCreated)
}

func (j *Job) SetStatus(status JobStatus) {
	j.Attributes[statusKey] = status
}

func (j *Job) recordCreationTime() {
	j.RecordTimeForAttribute(creationTimeKey)
}

func (j *Job) RecordCompletionTime() {
	j.RecordTimeForAttribute(completionTimeKey)
}

func (j *Job) RecordTimeForAttribute(key JobAttributeKey) {
	j.Attributes[key] = FormattedTimestamp()
}
