package api

import (
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/threading"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"
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
	muxUnderTest.Shutdown()
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

	muxUnderTest.Shutdown()
}

func TestPostScenarioTomlResource_OkResponse(t *testing.T) {
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
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	muxUnderTest.Shutdown()
}

func TestPostScenarioTextResource_BadRequestResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /scenario text request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: "This isn't TOML",
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	muxUnderTest.Shutdown()
}

func TestPostValidScenarioResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	scenarioTomlText := readTestFileAsText("testdata/ValidTestScenario.toml")

	// when
	postContext := TestContext{
		Name: "POST /scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: scenarioTomlText,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
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

	muxUnderTest.Shutdown()
}

func verifyResponseStatusCode(muxUnderTest *Mux, context TestContext) httptest.JsonResponseContainer {
	g := NewGomegaWithT(context.T)

	var responseContainer httptest.JsonResponseContainer

	responseContainer = sendRequest(muxUnderTest, context.Request)

	g.Expect(responseContainer.StatusCode).To(BeNumerically("==", context.ExpectedResponseStatus),
		context.Name+" should return status "+strconv.Itoa(context.ExpectedResponseStatus))

	return responseContainer
}

func buildMuxUnderTest() *Mux {
	threading.ResetMainThreadChannel()
	mainThreadChannel := threading.GetMainThreadChannel()
	muxUnderTest := new(Mux).Initialise().WithMainThreadChannel(&mainThreadChannel)
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

func verifyResponseTimeIsAboutNow(g *GomegaWithT, responseContainer httptest.JsonResponseContainer) {
	responseTimeString, ok := responseContainer.JsonMap["Time"].(string)
	g.Expect(ok).To(Equal(true), " should return a string encoding of time")

	responseTime, parseErr := time.Parse(time.RFC3339Nano, responseTimeString)
	g.Expect(parseErr).To(BeNil(), " should return a RFC3339Nano encoded string of time")

	g.Expect(responseTime).To(BeTemporally("~", time.Now(), time.Millisecond*5), " should return status time of about now")
}

func TestScenarioPutRequest_NotAllowedResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "PUT /scenario request returns 405 (not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PUT",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestInvalidModelPostScenario_BadRequestResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	scenarioTomlText := readTestFileAsText("testdata/InvalidModelTestScenario.toml")

	// when
	postContext := TestContext{
		Name: "POST /scenario text request with invalid model returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: scenarioTomlText,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	muxUnderTest.Shutdown()
}

func TestInvalidModelParameterPostScenario_BadRequestResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	scenarioTomlText := readTestFileAsText("testdata/InvalidModelParameterTestScenario.toml")

	// when
	postContext := TestContext{
		Name: "POST /scenario text request with invalid model parameter returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: scenarioTomlText,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	muxUnderTest.Shutdown()
}
