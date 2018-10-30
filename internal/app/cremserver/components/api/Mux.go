// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	serverApi "github.com/LindsayBradford/crem/server/api"
	"github.com/LindsayBradford/crem/server/job"
	"github.com/LindsayBradford/crem/server/rest"
)

const jobsPath = "jobs"

type JobArray []*job.Job

type Mux struct {
	serverApi.Mux

	jobs       *JobQueue
	JobHistory JobArray
}

func (m *Mux) Initialise() *Mux {
	m.Mux.Initialise()
	m.JobHistory = make(JobArray, 0)
	m.jobs = new(JobQueue).Initialise()
	m.jobs.JobFunction = m.DoJob

	m.AddHandler(BuildApiPath(jobsPath), m.V1HandleJobs)

	go m.jobs.Start()

	return m
}

func (m *Mux) WithJobQueueLength(length uint64) *Mux {
	m.jobs.WithQueueLength(length)
	return m
}

func (m *Mux) AddToHistory(newJob *job.Job) {
	m.JobHistory = append(JobArray{newJob}, m.JobHistory...) // for reverse chronological ordering
}

func (m *Mux) JobWithId(id string) *job.Job {
	var matchingJob *job.Job
	for _, historicalJob := range m.JobHistory {
		if string(historicalJob.Id) == id {
			matchingJob = historicalJob
		}
	}

	return matchingJob
}

func BuildApiPath(pathElements ...string) string {
	builtPath := baseApiPath()

	for _, element := range pathElements {
		builtPath = builtPath + rest.UrlPathSeparator + element
	}

	return builtPath
}

func baseApiPath() string {
	return serverApi.BasePath + serverApi.V1Path
}
