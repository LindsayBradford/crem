// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/config/interpreter"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment"
	serverApi "github.com/LindsayBradford/crem/internal/pkg/server/api"
	"github.com/LindsayBradford/crem/internal/pkg/server/job"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/threading"
	"io/ioutil"
	"net/http"
)

const v1Path = "v1"

const scenarioTextKey = "ScenarioText"
const scenarioNameKey = "ScenarioName"

type JobArray []*job.Job

type Mux struct {
	serverApi.Mux
	mainThreadChannel *threading.MainThreadChannel

	modelConfigInterpreter *interpreter.ModelConfigInterpreter
	model                  *catchment.Model
	modelSolution          *solution.Solution

	jsonMarshaler json.Marshaler

	attributes.ContainedAttributes
}

func (m *Mux) Initialise() *Mux {
	m.Mux.Initialise()

	m.modelConfigInterpreter = interpreter.NewModelConfigInterpreter()

	m.AddHandler(buildV1ApiPath("model"), m.v1modelHandler)
	m.AddHandler(buildV1ApiPath("scenario"), m.v1scenarioHandler)
	m.AddHandler(buildV1ApiPath("model/subcatchment/\\d+"), m.v1subcatchmentHandler)

	return m
}

func (m *Mux) WithMainThreadChannel(channel *threading.MainThreadChannel) *Mux {
	m.mainThreadChannel = channel
	return m
}

func (m *Mux) WithCacheMaxAge(maxAgeInSeconds uint64) *Mux {
	m.MuxImpl.WithCacheMaxAge(maxAgeInSeconds)
	return m
}

func buildV1ApiPath(pathElements ...string) string {
	builtPath := rest.UrlPathSeparator + serverApi.BasePath + rest.UrlPathSeparator + v1Path

	for _, element := range pathElements {
		builtPath = builtPath + rest.UrlPathSeparator + element
	}

	return "^" + builtPath + "$"
}

func (m *Mux) AddHandler(address string, handler rest.HandlerFunc) {
	m.HandlerMap.AddHandler(address, handler)
}

func (m *Mux) Shutdown() {
	if m.model != nil {
		m.model.TearDown()
	}
	m.MuxImpl.Shutdown()
}

func requestBodyToBytes(r *http.Request) []byte {
	responseBodyBytes, _ := ioutil.ReadAll(r.Body)
	return responseBodyBytes
}

func requestBodyToString(r *http.Request) string {
	return string(requestBodyToBytes(r))
}
