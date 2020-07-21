package api

import (
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

const (
	baseUrl = "http://dummyUrl/"
)

type TestContext struct {
	Name                   string
	T                      *testing.T
	Request                httptest.HttpTestRequestContext
	ExpectedResponseStatus int
}

func TestFirstScenarioGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "GET /scenario request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
}

func TestPostScenarioResource_NotAllowedResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /scenario text request returns 405 (not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: "here is some text that should be TOML",
			ContentType: rest.TextMimeType,
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then

	verifyResponseStatusCode(muxUnderTest, postContext)
}

func TestPostScenarioResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	scenarioTomlText := readTestFileAsText("testdata/ValidTestScenario.toml")

	// when
	postContext := TestContext{
		Name: "POST /scenario text request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: scenarioTomlText,
			ContentType: rest.TextMimeType,
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)
}

func TestScenarioResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	postContext := TestContext{
		Name: "POST /scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: "here is some text",
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then

	verifyResponseStatusCode(muxUnderTest, postContext)

	getContext := TestContext{
		Name: "GET /scenario request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
}

func verifyResponseStatusCode(muxUnderTest *Mux, context TestContext) {
	g := gomega.NewGomegaWithT(context.T)

	responseContainer := sendRequest(muxUnderTest, context.Request)

	g.Expect(responseContainer.StatusCode).To(gomega.BeNumerically("==", context.ExpectedResponseStatus),
		context.Name+" should return status "+strconv.Itoa(context.ExpectedResponseStatus))
}

func buildMuxUnderTest() *Mux {
	muxUnderTest := new(Mux).Initialise()
	muxUnderTest.SetLogger(loggers.DefaultTestingLogger)
	return muxUnderTest
}

func sendRequest(muxUnderTest *Mux, context httptest.HttpTestRequestContext) httptest.JsonResponseContainer {
	context.Handler = muxUnderTest.ServeHTTP
	return context.BuildJsonResponse()
}

func readTestFileAsText(filePath string) string {
	if b, err := ioutil.ReadFile(filePath); err == nil {
		return string(b)
	}
	return "error reading file"
}
