// Copyright (c) 2018 Australian Rivers Institute.

package test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
)

type HttpTestRequestContext struct {
	Method      string
	TargetUrl   string
	RequestBody string
	ContentType string
	Handler     http.HandlerFunc
}

type JsonResponseContainer struct {
	StatusCode  int
	JsonMap     map[string]interface{}
	RawResponse string
}

func (context *HttpTestRequestContext) BuildJsonResponse() JsonResponseContainer {
	response := context.buildResponse()

	responseAsJson, rawResponse := context.buildJsonMapFrom(response.Body)

	return JsonResponseContainer{
		StatusCode:  response.StatusCode,
		JsonMap:     responseAsJson,
		RawResponse: rawResponse,
	}
}

func (context *HttpTestRequestContext) buildResponse() *http.Response {
	w := httptest.NewRecorder()
	context.Handler(w, context.newRequest())
	return w.Result()
}

func (context *HttpTestRequestContext) newRequest() *http.Request {
	bodyReader := strings.NewReader(context.RequestBody)

	request := httptest.NewRequest(context.Method, context.TargetUrl, bodyReader)

	if context.ContentType != "" {
		request.Header.Add(rest.ContentTypeHeaderKey, context.ContentType)
	}
	return request
}

func (context *HttpTestRequestContext) buildJsonMapFrom(responseBody io.ReadCloser) (map[string]interface{}, string) {
	jsonMap := make(map[string]interface{})

	responseBodyBytes, _ := ioutil.ReadAll(responseBody)
	json.Unmarshal(responseBodyBytes, &jsonMap)

	return jsonMap, string(responseBodyBytes)
}
