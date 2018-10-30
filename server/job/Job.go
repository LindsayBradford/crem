// Copyright (c) 2018 Australian Rivers Institute.

package job

import (
	"github.com/LindsayBradford/crem/server"
	"github.com/LindsayBradford/crem/server/job/uuid"
)

const creationTimeKey = "CreationTime"
const completionTimeKey = "CompletedTime"

const statusKey = "Status"

type Id string
type Status AttributeKey
type AttributeKey string

const (
	Unspecified Status = "UNSPECIFIED"
	Created     Status = "CREATED"
	Completed   Status = "COMPLETED"
	Invalid     Status = "INVALID"
)

type Job struct {
	Id               Id
	Attributes       map[AttributeKey]interface{}
	HiddenAttributes map[AttributeKey]interface{} `json:"-"`
}

func (j *Job) Initialise() *Job {
	j.createNewJobID()
	j.makeAttributeMaps()
	j.recordCreationAttributes()
	return j
}

func (j *Job) createNewJobID() {
	j.Id = Id(uuid.New())
}

func (j *Job) makeAttributeMaps() {
	j.Attributes = make(map[AttributeKey]interface{}, 0)
	j.HiddenAttributes = make(map[AttributeKey]interface{}, 0)
}

func (j *Job) recordCreationAttributes() {
	j.recordCreationTime()
	j.SetStatus(Created)
}

func (j *Job) SetStatus(status Status) {
	j.Attributes[statusKey] = status
}

func (j *Job) Status() Status {
	status, ok := j.Attributes[statusKey].(Status)
	if ok {
		return status
	}
	return Unspecified
}

func (j *Job) IsProcessed() bool {
	if j.Status() == Completed || j.Status() == Invalid {
		return true
	}
	return false
}

func (j *Job) recordCreationTime() {
	j.RecordTimeForAttribute(creationTimeKey)
}

func (j *Job) RecordCompletionTime() {
	j.RecordTimeForAttribute(completionTimeKey)
}

func (j *Job) RecordTimeForAttribute(key AttributeKey) {
	j.Attributes[key] = server.FormattedTimestamp()
}
