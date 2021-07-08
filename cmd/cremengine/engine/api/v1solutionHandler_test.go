package api

import (
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	"net/http"
	"testing"
)

func TestSolutionsGetAsIsRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "GET /api/v1/solutions/As-Is request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/solutions/As-Is",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestGetSolutionsAsIsResource_OkResponse(t *testing.T) {
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
	postContext = TestContext{
		Name: "POST /solutions text request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/solutions",
			RequestBody: validSolutions,
			ContentType: rest.CsvMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	getContext := TestContext{
		Name: "GET /api/v1/solutions/As-Is request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/solutions/As-Is",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetSolutionsAThreeOfFourResource_OkResponse(t *testing.T) {
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
	postContext = TestContext{
		Name: "POST /solutions text request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/solutions",
			RequestBody: validSolutions,
			ContentType: rest.CsvMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	getContext := TestContext{
		Name: "GET /api/v1/solutions/3-of-4 request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    "GET",
			TargetUrl: baseUrl + "api/v1/solutions/3-of-8",
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetSolutionsMissingSolutionResource_OkResponse(t *testing.T) {
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
	postContext = TestContext{
		Name: "POST /solutions text request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/solutions",
			RequestBody: validSolutions,
			ContentType: rest.CsvMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	getContext := TestContext{
		Name: "GET /api/v1/solutions/MissingSolution request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    "GET",
			TargetUrl: baseUrl + "api/v1/solutions/MissingSolution",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestSolutionsPost_NotAllowedResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "POST /api/v1/solutions/testModel request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/solutions/testModel",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

//func TestGetModelsSPostScratchpad_OkResponse(t *testing.T) {
//	// given
//	muxUnderTest := buildMuxUnderTest()
//
//	// when
//	postContext := TestContext{
//		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
//		T:    t,
//		Request: httptest.HttpTestRequestContext{
//			Method:      "POST",
//			TargetUrl:   baseUrl + "api/v1/scenario",
//			RequestBody: validScenarioTomlConfig,
//			ContentType: rest.TomlMimeType,
//		},
//		ExpectedResponseStatus: http.StatusOK,
//	}
//
//	verifyResponseStatusCode(muxUnderTest, postContext)
//
//	// then
//
//	// when
//	requestAttributes := attributes.Attributes{
//		attributes.NameValuePair{
//			Name:  "Summary",
//			Value: "solution summary",
//		},
//		attributes.NameValuePair{
//			Name:  "Encoding",
//			Value: "A1",
//		},
//	}
//
//	attributesAsJson, _ := json.Marshal(requestAttributes)
//	attributesAsJsonString := string(attributesAsJson)
//
//	getContext := TestContext{
//		Name: "POST /api/v1/models/Scratchpad request returns 200 (ok) response",
//		T:    t,
//		Request: httptest.HttpTestRequestContext{
//			Method:      "POST",
//			TargetUrl:   baseUrl + "api/v1/models/Scratchpad",
//			ContentType: rest.JsonMimeType,
//			RequestBody: attributesAsJsonString,
//		},
//		ExpectedResponseStatus: http.StatusOK,
//	}
//
//	// then
//	verifyResponseStatusCode(muxUnderTest, getContext)
//
//	// when
//	requestAttributes = attributes.Attributes{
//		attributes.NameValuePair{
//			Name:  "Encoding",
//			Value: "A3",
//		},
//	}
//
//	attributesAsJson, _ = json.Marshal(requestAttributes)
//	attributesAsJsonString = string(attributesAsJson)
//
//	getContext = TestContext{
//		Name: "POST /api/v1/models/Scratchpad request returns 200 (ok) response",
//		T:    t,
//		Request: httptest.HttpTestRequestContext{
//			Method:      "POST",
//			TargetUrl:   baseUrl + "api/v1/models/Scratchpad",
//			ContentType: rest.JsonMimeType,
//			RequestBody: attributesAsJsonString,
//		},
//		ExpectedResponseStatus: http.StatusOK,
//	}
//
//	// then
//	verifyResponseStatusCode(muxUnderTest, getContext)
//
//	muxUnderTest.Shutdown()
//}
