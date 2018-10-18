// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crm/server"
	"github.com/nu7hatch/gouuid"
)

const jobsPath = "/jobs"

type CrmApiMux struct {
	server.ApiMux

	jobs *JobQueue
}

func (cam *CrmApiMux) Initialise() *CrmApiMux {
	cam.ApiMux.Initialise()
	cam.jobs = new(JobQueue).Initialise()
	cam.AddHandler(baseApiPath()+jobsPath, cam.V1HandleJobs)
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
	dequeuedJob, updatedQueue := jq.Jobs[len(jq.Jobs)-1], jq.Jobs[:len(jq.Jobs)-1]
	jq.Jobs = updatedQueue
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

func (ji *Job) Initialise() *Job {
	newId, _ := uuid.NewV4()
	ji.JobId = newId.String()

	ji.Attributes = make(map[string]interface{}, 0)
	ji.HiddenAttributes = make(map[string]interface{}, 0)

	ji.Attributes["CreationTime"] = server.FormattedTimestamp()
	ji.Attributes["Status"] = "CREATED"
	return ji
}
