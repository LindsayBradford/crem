package api

import (
	"encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"net/http"
	"testing"
)

func TestModelsGetAsIsRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "GET /api/v1/models/As-Is request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/models/As-Is",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModelsGetScratchpadRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "GET /api/v1/models/Scratchpad request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/models/Scratchpad",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestGetModelsAsIsResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	verifyResponseStatusCode(muxUnderTest, postContext)

	// then

	// when
	getContext := TestContext{
		Name: "GET /api/v1/models/As-Is request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/models/As-Is",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetModelsScratchpadResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	verifyResponseStatusCode(muxUnderTest, postContext)

	// then

	// when
	getContext := TestContext{
		Name: "GET /api/v1/models/Scratchpad request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/models/Scratchpad",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestModelsPostBeforeScenario_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "POST /api/v1/models/testModel request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/models/testModel",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestGetModelsSPostAsIsModel_NotAllowedResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	verifyResponseStatusCode(muxUnderTest, postContext)

	// then

	// when
	getContext := TestContext{
		Name: "POST /api/v1/models/As-Is request returns 405 (not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/models/As-Is",
			ContentType: rest.JsonMimeType,
			RequestBody: "{}",
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetModelsSPostNotJsonMimeType_UnsupportedMediaTypeResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	verifyResponseStatusCode(muxUnderTest, postContext)

	// then

	// when
	getContext := TestContext{
		Name: "POST /api/v1/models/testModel request returns 415 (unsupported media type) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/models/testModel",
			ContentType: rest.TomlMimeType,
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusUnsupportedMediaType,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetModelsSPostNonJsonBody_BadRequestResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	verifyResponseStatusCode(muxUnderTest, postContext)

	// then

	// when
	getContext := TestContext{
		Name: "POST /api/v1/models/testModel request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/models/testModel",
			ContentType: rest.JsonMimeType,
			RequestBody: "this is not json content",
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetModelsSPostScratchpad_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	verifyResponseStatusCode(muxUnderTest, postContext)

	// then

	// when
	requestAttributes := attributes.Attributes{
		attributes.NameValuePair{
			Name:  "Summary",
			Value: "solution summary",
		},
		attributes.NameValuePair{
			Name:  "Encoding",
			Value: "A1",
		},
	}

	attributesAsJson, _ := json.Marshal(requestAttributes)
	attributesAsJsonString := string(attributesAsJson)

	getContext := TestContext{
		Name: "POST /api/v1/models/Scratchpad request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/models/Scratchpad",
			ContentType: rest.JsonMimeType,
			RequestBody: attributesAsJsonString,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)

	// when
	requestAttributes = attributes.Attributes{
		attributes.NameValuePair{
			Name:  "Encoding",
			Value: "A3",
		},
	}

	attributesAsJson, _ = json.Marshal(requestAttributes)
	attributesAsJsonString = string(attributesAsJson)

	getContext = TestContext{
		Name: "POST /api/v1/models/Scratchpad request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/models/Scratchpad",
			ContentType: rest.JsonMimeType,
			RequestBody: attributesAsJsonString,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)

	muxUnderTest.Shutdown()
}

func TestGetModelsPostMissing_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	getContext := TestContext{
		Name: "GET /api/v1/models/MissingModel request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    "GET",
			TargetUrl: baseUrl + "api/v1/models/MissingModel",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetModelsPostNew_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	verifyResponseStatusCode(muxUnderTest, postContext)

	// then

	// when
	requestAttributes := attributes.Attributes{
		attributes.NameValuePair{
			Name:  "Summary",
			Value: "solution summary",
		},
		attributes.NameValuePair{
			Name:  "Encoding",
			Value: "A1",
		},
	}

	attributesAsJson, _ := json.Marshal(requestAttributes)
	attributesAsJsonString := string(attributesAsJson)

	getContext := TestContext{
		Name: "POST /api/v1/models/NewModel request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/models/NewModel",
			ContentType: rest.JsonMimeType,
			RequestBody: attributesAsJsonString,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetModelsTwicePostNew_NotAllowedResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	requestAttributes := attributes.Attributes{
		attributes.NameValuePair{
			Name:  "Summary",
			Value: "solution summary",
		},
		attributes.NameValuePair{
			Name:  "Encoding",
			Value: "A1",
		},
	}

	attributesAsJson, _ := json.Marshal(requestAttributes)
	attributesAsJsonString := string(attributesAsJson)

	getContext := TestContext{
		Name: "POST /api/v1/models/NewModel request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/models/NewModel",
			ContentType: rest.JsonMimeType,
			RequestBody: attributesAsJsonString,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// when
	requestAttributes = attributes.Attributes{
		attributes.NameValuePair{
			Name:  "Summary",
			Value: "solution summary2",
		},
		attributes.NameValuePair{
			Name:  "Encoding",
			Value: "A0",
		},
	}
	// then
	verifyResponseStatusCode(muxUnderTest, getContext)

	// when
	attributesAsJson, _ = json.Marshal(requestAttributes)
	attributesAsJsonString = string(attributesAsJson)

	getContext = TestContext{
		Name: "POST /api/v1/models/NewModel request returns 405 (method not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/models/NewModel",
			ContentType: rest.JsonMimeType,
			RequestBody: attributesAsJsonString,
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}
