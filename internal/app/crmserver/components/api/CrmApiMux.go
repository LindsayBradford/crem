// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crm/server"
	"github.com/nu7hatch/gouuid"
)

const jobsPath = "jobs"

const creationTimeKey = "CreationTime"
const statusKey = "Status"

const defaultQueueLength = 1 // TODO: Make this configurable.

type CrmApiMux struct {
	server.ApiMux

	jobs       *CrmApiJobQueue
	JobHistory []*Job
}

func (cam *CrmApiMux) Initialise() *CrmApiMux {
	cam.ApiMux.Initialise()
	cam.JobHistory = make([]*Job, 0)
	cam.jobs = new(CrmApiJobQueue).Initialise()
	cam.jobs.JobFunction = cam.DoJob

	cam.AddHandler(BuildApiPath(jobsPath), cam.V1HandleJobs)

	go cam.jobs.Start()

	return cam
}

func (cam *CrmApiMux) AddToHistory(newJob *Job) {
	cam.JobHistory = append([]*Job{newJob}, cam.JobHistory...) // for reverse chronological ordering

}

func (cam *CrmApiMux) JobWithId(id string) *Job {
	var matchingJob *Job
	for _, job := range cam.JobHistory {
		if job.JobId == id {
			matchingJob = job
		}
	}

	return matchingJob
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

type JobQueue struct { // TODO: Split this out into its own separate generic Go file.
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

type CrmApiJobQueue struct { // TODO: Put convenience methods around this to make it easier to use for CrmApi.
	JobQueue
}

func (q *CrmApiJobQueue) Initialise() *CrmApiJobQueue {
	q.JobQueue.Initialise()
	return q
}
