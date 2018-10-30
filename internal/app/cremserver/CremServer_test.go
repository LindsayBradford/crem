// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/LindsayBradford/crem/internal/app/cremserver/components"
	"github.com/LindsayBradford/crem/internal/app/cremserver/components/api"
	"github.com/LindsayBradford/crem/logging/handlers"
	"github.com/LindsayBradford/crem/server/rest"
	"github.com/LindsayBradford/crem/server/test"
	. "github.com/onsi/gomega"
)

type testContext struct {
	name       string
	t          *testing.T
	configFile string
}

func TestDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := testContext{
		name:       "Single run of Dumb annealer",
		t:          t,
		configFile: "testdata/DumbAnnealerTestConfig-OneRun.toml",
	}

	verifyDumbAnnealerRunsAgainstContext(context)
}

func TestDumbAnnealerIntegrationThreeRunsSequentially(t *testing.T) {
	context := testContext{
		name:       "Three sequential runs of Dumb annealer",
		t:          t,
		configFile: "testdata/DumbAnnealerTestConfig-ThreeRunsSequentially.toml",
	}

	verifyDumbAnnealerRunsAgainstContext(context)
}

func TestDumbAnnealerIntegrationThreeRunsConcurrently(t *testing.T) {
	context := testContext{
		name:       "Three concurrent runs of Dumb annealer",
		t:          t,
		configFile: "testdata/DumbAnnealerTestConfig-ThreeRunsConcurrently.toml",
	}

	verifyDumbAnnealerRunsAgainstContext(context)
}

func verifyDumbAnnealerRunsAgainstContext(context testContext) {
	if testing.Short() {
		context.t.Skip("skipping " + context.name + " in short mode")
	}
	g := NewGomegaWithT(context.t)

	simulatedMainCall := func() {
		components.RunScenarioFromConfigFile(context.configFile)
	}

	g.Expect(simulatedMainCall).To(Not(Panic()), context.name+" should not panic")
}

func TestValidJobsGetRequest_OkResponse(t *testing.T) {
	context := testContext{
		name:       "GET /jobs request returns 200 response",
		t:          t,
		configFile: "testdata/server.toml",
	}

	verifyResponseToValidJobsGetRequest(context)
}

func verifyResponseToValidJobsGetRequest(context testContext) {

	g := NewGomegaWithT(context.t)

	muxUnderTest := buildMuxUnderTest()

	requestContext := test.HttpTestRequestContext{
		Method:    "GET",
		TargetUrl: "http://dummyUrl/api/v1/jobs",
		Handler:   muxUnderTest.V1HandleJobs,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(BeNumerically("==", http.StatusOK), context.name+" should return OK status")

	// verifyResponseTimeIsAboutNow(g, responseContainer)
}

func TestInvalidValidJobsPostRequest_InternalServerErrorResponse(t *testing.T) {
	context := testContext{
		name:       "POST /jobs request of invalid scenario returns 500 response",
		t:          t,
		configFile: "testdata/server.toml",
	}

	verifyInternalServerErrorResponseToInvalidJobsPostRequest(context)
}

func verifyInternalServerErrorResponseToInvalidJobsPostRequest(context testContext) {

	g := NewGomegaWithT(context.t)

	muxUnderTest := buildMuxUnderTest()

	requestContext := test.HttpTestRequestContext{
		Method:      "POST",
		TargetUrl:   "http://dummyUrl/api/v1/jobs",
		ContentType: rest.TomlMimeType,
		RequestBody: "invalidScenarioText: isInvalid",
		Handler:     muxUnderTest.V1HandleJobs,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(BeNumerically("==", http.StatusInternalServerError), context.name+" should return 500 status")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func buildMuxUnderTest() *api.Mux {
	muxUnderTest := new(api.Mux).Initialise()
	muxUnderTest.SetLogger(handlers.DefaultTestingLogHandler)
	return muxUnderTest
}

func verifyResponseTimeIsAboutNow(g *GomegaWithT, responseContainer test.JsonResponseContainer) {
	responseTimeString, ok := responseContainer.JsonMap["Time"].(string)
	g.Expect(ok).To(Equal(true), " should return a string encoding of time")

	responseTime, parseErr := time.Parse(time.RFC3339Nano, responseTimeString)
	g.Expect(parseErr).To(BeNil(), " should return a RFC3339Nano encoded string of time")

	g.Expect(responseTime).To(BeTemporally("~", time.Now(), time.Millisecond*5), " should return status time of about now")
}
