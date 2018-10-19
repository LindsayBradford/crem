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

	jobs *JobQueue
}

func (cam *CrmApiMux) Initialise() *CrmApiMux {
	cam.ApiMux.Initialise()
	cam.jobs = new(JobQueue).Initialise()
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
	JobHistory  []Job
	Jobs2       chan Job      `json:"-"`
	JobFunction func(job Job) `json:"-"`
}

func (jq *JobQueue) Initialise() *JobQueue {
	jq.JobHistory = make([]Job, 0)
	jq.Jobs2 = make(chan Job, defaultQueueLength)
	return jq
}

const enqueueSucceeded = true
const enqueueFailed = false

func (jq *JobQueue) Enqueue(newJob Job) bool {
	jq.JobHistory = append(jq.JobHistory, newJob)
	select {
	case jq.Jobs2 <- newJob:
		return enqueueSucceeded
	default:
		return enqueueFailed
	}
}

func (jq *JobQueue) Start() {
	for {
		job := <-jq.Jobs2
		jq.JobFunction(job)
	}
}

func (jq *JobQueue) JobWithId(id string) Job {
	var matchingJob Job
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
