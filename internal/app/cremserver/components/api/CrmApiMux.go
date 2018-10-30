// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crem/server"
	"github.com/LindsayBradford/crem/server/api"
	"github.com/LindsayBradford/crem/server/job"
)

const jobsPath = "jobs"

type CremApiMux struct {
	api.Mux

	jobs       *CremApiJobQueue
	JobHistory []*job.Job
}

func (cam *CremApiMux) Initialise() *CremApiMux {
	cam.Mux.Initialise()
	cam.JobHistory = make([]*job.Job, 0)
	cam.jobs = new(CremApiJobQueue).Initialise()
	cam.jobs.JobFunction = cam.DoJob

	cam.AddHandler(BuildApiPath(jobsPath), cam.V1HandleJobs)

	go cam.jobs.Start()

	return cam
}

func (cam *CremApiMux) WithJobQueueLength(length uint64) *CremApiMux {
	cam.jobs.WithQueueLength(length)
	return cam
}

func (cam *CremApiMux) AddToHistory(newJob *job.Job) {
	cam.JobHistory = append([]*job.Job{newJob}, cam.JobHistory...) // for reverse chronological ordering
}

func (cam *CremApiMux) JobWithId(id string) *job.Job {
	idSupplied := job.Id(id)

	var matchingJob *job.Job
	for _, historicalJob := range cam.JobHistory {
		if historicalJob.Id == idSupplied {
			matchingJob = historicalJob
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
	return api.BasePath + api.V1Path
}

type CremApiJobQueue struct {
	job.Queue
}

func (q *CremApiJobQueue) Initialise() *CremApiJobQueue {
	q.Queue.Initialise()
	return q
}

func (q *CremApiJobQueue) WithQueueLength(length uint64) *CremApiJobQueue {
	q.Queue.WithQueueLength(length)
	return q
}
