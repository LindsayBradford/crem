// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	serverApi "github.com/LindsayBradford/crem/internal/pkg/server/api"
	"github.com/LindsayBradford/crem/internal/pkg/server/job"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/threading"
)

const v1Path = "v1"

const scenarioTextKey = "ScenarioText"

type JobArray []*job.Job

type Mux struct {
	serverApi.Mux
	mainThreadChannel *threading.MainThreadChannel

	attributes.ContainedAttributes
}

func (m *Mux) Initialise() *Mux {
	m.Mux.Initialise()

	m.AddHandler(buildV1ApiPath("scenario"), m.v1scenarioHandler)

	return m
}

func (m *Mux) WithMainThreadChannel(channel *threading.MainThreadChannel) *Mux {
	m.mainThreadChannel = channel
	return m
}

func buildV1ApiPath(pathElements ...string) string {
	builtPath := rest.UrlPathSeparator + serverApi.BasePath + rest.UrlPathSeparator + v1Path

	for _, element := range pathElements {
		builtPath = builtPath + rest.UrlPathSeparator + element
	}

	return builtPath
}
