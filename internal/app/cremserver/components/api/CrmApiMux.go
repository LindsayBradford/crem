// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crem/server"
)

const jobsPath = "jobs"

type CremApiMux struct {
	server.ApiMux

	jobs       *CremApiJobQueue
	JobHistory []*server.Job
}

func (cam *CremApiMux) Initialise() *CremApiMux {
	cam.ApiMux.Initialise()
	cam.JobHistory = make([]*server.Job, 0)
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

func (cam *CremApiMux) AddToHistory(newJob *server.Job) {
	cam.JobHistory = append([]*server.Job{newJob}, cam.JobHistory...) // for reverse chronological ordering
}

func (cam *CremApiMux) JobWithId(id string) *server.Job {
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

type CremApiJobQueue struct {
	server.JobQueue
}

func (q *CremApiJobQueue) Initialise() *CremApiJobQueue {
	q.JobQueue.Initialise()
	return q
}

func (q *CremApiJobQueue) WithQueueLength(length uint64) *CremApiJobQueue {
	q.JobQueue.WithQueueLength(length)
	return q
}
