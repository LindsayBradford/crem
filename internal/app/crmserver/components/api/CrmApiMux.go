// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crm/server"
)

const jobsPath = "jobs"

type CrmApiMux struct {
	server.ApiMux

	jobs       *CrmApiJobQueue
	JobHistory []*server.Job
}

func (cam *CrmApiMux) Initialise() *CrmApiMux {
	cam.ApiMux.Initialise()
	cam.JobHistory = make([]*server.Job, 0)
	cam.jobs = new(CrmApiJobQueue).Initialise()
	cam.jobs.JobFunction = cam.DoJob

	cam.AddHandler(BuildApiPath(jobsPath), cam.V1HandleJobs)

	go cam.jobs.Start()

	return cam
}

func (cam *CrmApiMux) AddToHistory(newJob *server.Job) {
	cam.JobHistory = append([]*server.Job{newJob}, cam.JobHistory...) // for reverse chronological ordering
}

func (cam *CrmApiMux) JobWithId(id string) *server.Job {
	idSupplied := server.JobId(id)

	var matchingJob *server.Job
	for _, job := range cam.JobHistory {
		if job.Id == idSupplied {
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

type CrmApiJobQueue struct {
	server.JobQueue
}

func (q *CrmApiJobQueue) Initialise() *CrmApiJobQueue {
	q.JobQueue.Initialise()
	return q
}
