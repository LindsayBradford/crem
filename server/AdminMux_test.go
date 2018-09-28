// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/onsi/gomega"
)

type testContext struct {
	name       string
	t          *testing.T
	configFile string
}

type RequestContext struct {
	method    string
	targetUrl string
	handler   http.HandlerFunc
}

func TestValidStatusRequest_OkResponse(t *testing.T) {
	context := testContext{
		name:       "GET /status request returns 200 response",
		t:          t,
		configFile: "testdata/server.toml",
	}

	verifyResponseToValidStatusRequest(context)
}

func TestInvalidStatusRequest_MethodNotAllowedResponse(t *testing.T) {
	context := testContext{
		name:       "POST /status request returns 405 response",
		t:          t,
		configFile: "testdata/server.toml",
	}

	verifyResponseToInvalidStatusRequest(context)
}

type ResponseContainer struct {
	statusCode int
	jsonBody   map[string]interface{}
}

func verifyResponseToValidStatusRequest(context testContext) {
	g := NewGomegaWithT(context.t)

	muxUnderTest := buildMuxUnderTest()

	expectedMessage := "Some bogus status message"
	expectedName := "Some bogus name"
	expectedVersion := "some bogus version"
	bogusTime := "some bogus time"
	muxUnderTest.Status = Status{Name: expectedName, Version: expectedVersion, Message: expectedMessage, Time: bogusTime}

	requestContext := RequestContext{
		method:    "GET",
		targetUrl: "http://dummyUrl/status",
		handler:   muxUnderTest.statusHandler,
	}

	responseContainer := parseResponse(getResponseFor(requestContext))

	g.Expect(responseContainer.statusCode).To(BeNumerically("==", http.StatusOK), context.name+" should return OK status")
	g.Expect(responseContainer.jsonBody["Name"]).To(Equal(expectedName), context.name+" should return expected status name")
	g.Expect(responseContainer.jsonBody["Version"]).To(Equal(expectedVersion), context.name+" should return expected status version")
	g.Expect(responseContainer.jsonBody["Message"]).To(Equal(expectedMessage), context.name+" should return expected status message")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func verifyResponseToInvalidStatusRequest(context testContext) {
	g := NewGomegaWithT(context.t)

	muxUnderTest := buildMuxUnderTest()

	requestContext := RequestContext{
		method:    "POST",
		targetUrl: "http://dummyUrl/status",
		handler:   muxUnderTest.statusHandler,
	}

	responseContainer := parseResponse(getResponseFor(requestContext))

	expectedResponseCode := http.StatusMethodNotAllowed
	g.Expect(responseContainer.statusCode).To(BeNumerically("==", expectedResponseCode), context.name+" should return Method not Allowed status")
	g.Expect(responseContainer.jsonBody["ResponseCode"]).To(BeNumerically("==", expectedResponseCode), context.name+" should return expected status code")

	expectedMessage := "Method not allowed"
	g.Expect(responseContainer.jsonBody["Message"]).To(Equal(expectedMessage), context.name+" should return expected status message")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func buildMuxUnderTest() *AdminMux {
	muxUnderTest := new(AdminMux).Initialise()
	muxUnderTest.SetLogger(handlers.DefaultNullLogHandler)
	return muxUnderTest
}

func getResponseFor(context RequestContext) *http.Response {
	req := httptest.NewRequest(context.method, context.targetUrl, nil)
	w := httptest.NewRecorder()

	muxUnderTest := new(AdminMux).Initialise()
	muxUnderTest.SetLogger(handlers.DefaultNullLogHandler)

	context.handler(w, req)

	return w.Result()
}

func parseResponse(response *http.Response) *ResponseContainer {
	container := new(ResponseContainer)
	container.statusCode = response.StatusCode

	responseBodyBytes, _ := ioutil.ReadAll(response.Body)

	container.jsonBody = make(map[string]interface{})
	json.Unmarshal(responseBodyBytes, &container.jsonBody)

	return container
}

func verifyResponseTimeIsAboutNow(g *GomegaWithT, responseContainer *ResponseContainer) {
	responseTimeString, ok := responseContainer.jsonBody["Time"].(string)
	g.Expect(ok).To(Equal(true), " should return a string encoding of time")

	responseTime, parseErr := time.Parse(time.RFC3339Nano, responseTimeString)
	g.Expect(parseErr).To(BeNil(), " should return a RFC3339Nano encoded string of time")

	g.Expect(responseTime).To(BeTemporally("~", time.Now(), time.Millisecond*5), " should return status time of about now")
}

func TestValidShutdownRequest_OkResponse(t *testing.T) {
	context := testContext{
		name:       "POST /shutdown request returns 200 response",
		t:          t,
		configFile: "testdata/server.toml",
	}

	verifyResponseToValidShutdownRequest(context)
}

func verifyResponseToValidShutdownRequest(context testContext) {
	g := NewGomegaWithT(context.t)

	muxUnderTest := buildMuxUnderTest()

	requestContext := RequestContext{
		method:    "POST",
		targetUrl: "http://dummyUrl/shutdown",
		handler:   muxUnderTest.shutdownHandler,
	}

	var responseContainer *ResponseContainer
	go func() {
		responseContainer = parseResponse(getResponseFor(requestContext))
	}()

	waitFinished := true
	waitFunc := func() bool {
		muxUnderTest.WaitForShutdownSignal()
		return true
	}()

	g.Eventually(waitFunc).Should(Equal(waitFinished), context.name+" should finish Wait function")
	g.Expect(responseContainer.statusCode).To(BeNumerically("==", http.StatusOK), context.name+" should return OK status")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func TestInvalidShutdownRequest_MethodNotAllowedResponse(t *testing.T) {
	context := testContext{
		name:       "GET /shutdown request returns 405 response",
		t:          t,
		configFile: "testdata/server.toml",
	}

	verifyResponseToInvalidShutdownRequest(context)
}

func verifyResponseToInvalidShutdownRequest(context testContext) {
	g := NewGomegaWithT(context.t)

	muxUnderTest := buildMuxUnderTest()

	requestContext := RequestContext{
		method:    "GET",
		targetUrl: "http://dummyUrl/shutdown",
		handler:   muxUnderTest.shutdownHandler,
	}

	var responseContainer *ResponseContainer
	responseContainer = parseResponse(getResponseFor(requestContext))

	expectedResponseCode := http.StatusMethodNotAllowed
	g.Expect(responseContainer.statusCode).To(BeNumerically("==", expectedResponseCode), context.name+" should return Method not Allowed status")
	g.Expect(responseContainer.jsonBody["ResponseCode"]).To(BeNumerically("==", expectedResponseCode), context.name+" should return expected status code")

	expectedMessage := "Method not allowed"
	g.Expect(responseContainer.jsonBody["Message"]).To(Equal(expectedMessage), context.name+" should return expected status message")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}
