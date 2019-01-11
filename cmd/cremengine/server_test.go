// Copyright (c) 2019 Australian Rivers Institute.

package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/LindsayBradford/crem/cmd/cremengine/components/api"
	testing2 "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/internal/pkg/server/test"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/onsi/gomega"
)

func TestValidJobsGetRequest_OkResponse(t *testing.T) {
	context := testing2.Context{
		Name:           "GET /jobs request returns 200 response",
		T:              t,
		ConfigFilePath: "testdata/server.toml",
	}

	verifyResponseToValidJobsGetRequest(context)
}

func verifyResponseToValidJobsGetRequest(context testing2.Context) {

	g := gomega.NewGomegaWithT(context.T)

	muxUnderTest := buildMuxUnderTest()

	requestContext := test.HttpTestRequestContext{
		Method:    "GET",
		TargetUrl: "http://dummyUrl/api/v1/jobs",
		Handler:   muxUnderTest.ServeHTTP,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(gomega.BeNumerically("==", http.StatusOK), context.Name+" should return OK status")
}

func TestInvalidValidJobsPostRequest_InternalServerErrorResponse(t *testing.T) {
	context := testing2.Context{
		Name:           "POST /jobs request of invalid scenario returns 500 response",
		T:              t,
		ConfigFilePath: "testdata/server.toml",
	}

	verifyInternalServerErrorResponseToInvalidJobsPostRequest(context)
}

func verifyInternalServerErrorResponseToInvalidJobsPostRequest(context testing2.Context) {

	g := gomega.NewGomegaWithT(context.T)

	muxUnderTest := buildMuxUnderTest()

	requestContext := test.HttpTestRequestContext{
		Method:      "POST",
		TargetUrl:   "http://dummyUrl/api/v1/jobs",
		ContentType: rest.TomlMimeType,
		RequestBody: "invalidScenarioText: isInvalid",
		Handler:     muxUnderTest.ServeHTTP,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(gomega.BeNumerically("==", http.StatusInternalServerError), context.Name+" should return 500 status")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func buildMuxUnderTest() *api.Mux {
	muxUnderTest := new(api.Mux).Initialise()
	muxUnderTest.SetLogger(loggers.DefaultTestingLogger)
	return muxUnderTest
}

func verifyResponseTimeIsAboutNow(g *gomega.GomegaWithT, responseContainer test.JsonResponseContainer) {
	responseTimeString, ok := responseContainer.JsonMap["Time"].(string)
	g.Expect(ok).To(gomega.Equal(true), " should return a string encoding of time")

	responseTime, parseErr := time.Parse(time.RFC3339Nano, responseTimeString)
	g.Expect(parseErr).To(gomega.BeNil(), " should return a RFC3339Nano encoded string of time")

	g.Expect(responseTime).To(gomega.BeTemporally("~", time.Now(), time.Millisecond*5), " should return status time of about now")
}
