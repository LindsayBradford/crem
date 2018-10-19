// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crm/server"
	"github.com/nu7hatch/gouuid"
)

const jobsPath = "jobs"

const creationTimeKey = "CreationTime"
const statusKey = "Status"

const defaultQueueLength = 1

type CrmApiMux struct {
	server.ApiMux

	jobs *CrmApiJobQueue
}

func (cam *CrmApiMux) Initialise() *CrmApiMux {
	cam.ApiMux.Initialise()
	cam.jobs = new(CrmApiJobQueue).Initialise()
	cam.jobs.JobFunction = cam.DoJob

	cam.AddHandler(BuildApiPath(jobsPath), cam.V1HandleJobs)

	go cam.jobs.Start()

	return cam
}

func BuildApiPath(pathElements ...string) string {
	builtPath := baseApiPath()

	for _, element := range pathElements {
		builtPath = builtPath + server.UrlPathSeparator + element
	}

	return builtPath
}

func baseApiPath() string {
	return server.ApiPath + server.V1Path
}

type JobQueue struct {
	JobHistory  []*Job
	Jobs        chan *Job      `json:"-"`
	JobFunction func(job *Job) `json:"-"`
}

func (jq *JobQueue) Initialise() *JobQueue {
	jq.JobHistory = make([]*Job, 0)
	jq.Jobs = make(chan *Job, defaultQueueLength)
	return jq
}

type JobEnqueuedStatus bool

const JobEnqueueSucceeded JobEnqueuedStatus = true
const JobEnqueueFailed JobEnqueuedStatus = false

func (jq *JobQueue) Enqueue(newJob *Job) JobEnqueuedStatus {
	jq.AddToHistory(newJob)
	select {
	case jq.Jobs <- newJob:
		return JobEnqueueSucceeded
	default:
		return JobEnqueueFailed
	}
}

func (jq *JobQueue) AddToHistory(newJob *Job) {
	jq.JobHistory = append([]*Job{newJob}, jq.JobHistory...) // for reverse chronological ordering

}

func (jq *JobQueue) Start() {
	for {
		job := <-jq.Jobs
		jq.JobFunction(job)
	}
}

func (jq *JobQueue) JobWithId(id string) *Job {
	var matchingJob *Job
	for _, job := range jq.JobHistory {
		if job.JobId == id {
			matchingJob = job
		}
	}

	return matchingJob
}

type Job struct {
	JobId            string
	Attributes       map[string]interface{}
	HiddenAttributes map[string]interface{} `json:"-"`
}

func (j *Job) Initialise() *Job {
	j.createNewJobID()
	j.makeAttributeMaps()
	j.recordCreationAttributes()
	return j
}

func (j *Job) createNewJobID() {
	newId, _ := uuid.NewV4()
	j.JobId = newId.String()
}

func (j *Job) makeAttributeMaps() {
	j.Attributes = make(map[string]interface{}, 0)
	j.HiddenAttributes = make(map[string]interface{}, 0)
}

func (j *Job) recordCreationAttributes() {
	j.Attributes[creationTimeKey] = server.FormattedTimestamp()
	j.Attributes[statusKey] = "CREATED"
}

type CrmApiJobQueue struct {
	JobQueue
}

func (q *CrmApiJobQueue) Initialise() *CrmApiJobQueue {
	q.JobQueue.Initialise()
	return q
}
