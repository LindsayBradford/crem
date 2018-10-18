// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crm/server"
	"github.com/nu7hatch/gouuid"
)

const jobsPath = "jobs"

const creationTimeKey = "CreationTime"
const statusKey = "Status"

type CrmApiMux struct {
	server.ApiMux

	jobs *JobQueue
}

func (cam *CrmApiMux) Initialise() *CrmApiMux {
	cam.ApiMux.Initialise()
	cam.jobs = new(JobQueue).Initialise()
	cam.AddHandler(BuildApiPath(jobsPath), cam.V1HandleJobs)
	return cam
}

func baseApiPath() string {
	return server.ApiPath + server.V1Path
}

type JobQueue struct {
	Jobs []Job
}

func (jq *JobQueue) Initialise() *JobQueue {
	jq.Jobs = make([]Job, 0)
	return jq
}

func (jq *JobQueue) Enqueue(newJob Job) {
	jq.Jobs = append(jq.Jobs, newJob)
}

func (jq *JobQueue) Dequeue() Job {
	dequeuedJob, queueSansJob := jq.Jobs[len(jq.Jobs)-1], jq.Jobs[:len(jq.Jobs)-1]
	jq.Jobs = queueSansJob
	return dequeuedJob
}

func (jq *JobQueue) JobWithId(id string) Job {
	var matchingJob Job
	for _, job := range jq.Jobs {
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

func BuildApiPath(pathElements ...string) string {
	builtPath := baseApiPath()

	for _, element := range pathElements {
		builtPath = builtPath + server.UrlPathSeparator + element
	}

	return builtPath
}
