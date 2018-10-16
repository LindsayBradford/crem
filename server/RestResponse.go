// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crm/errors"
	errors2 "github.com/pkg/errors"
)

type RestResponse struct {
	ResponseCode int
	ContentType  string
	CacheControl string
	Content      string
	Writer       http.ResponseWriter

	errors *errors.CompositeError
}

func (rrc *RestResponse) Initialise() *RestResponse {
	rrc.errors = new(errors.CompositeError)
	rrc.WithContentType(DefaultResposneContentType)
	return rrc
}

func (rrc *RestResponse) WithWriter(writer http.ResponseWriter) *RestResponse {
	rrc.Writer = writer
	return rrc
}

func (rrc *RestResponse) WithResponseCode(responseCode int) *RestResponse {
	rrc.ResponseCode = responseCode
	return rrc
}

func (rrc *RestResponse) WithCacheControlMaxAge(cacheMaxAgeInSeconds uint64) *RestResponse {
	rrc.CacheControl = fmt.Sprintf("max-age=%d", cacheMaxAgeInSeconds)
	return rrc
}

func (rrc *RestResponse) WithContentType(contentType string) *RestResponse {
	rrc.ContentType = contentType
	return rrc
}

func (rrc *RestResponse) WithJsonContent(content interface{}) *RestResponse {
	rrc.WithContentType(JsonMimeType)

	contentAsJsonBytes, encodeError := json.MarshalIndent(content, "", "  ")
	if encodeError != nil {
		wrappingError := errors2.Wrap(encodeError, "json content encoding")
		rrc.errors.Add(wrappingError)
	} else {
		rrc.Content = string(contentAsJsonBytes)
	}
	return rrc
}

func (rrc *RestResponse) Write() error {
	rrc.writeHeader()
	rrc.writeBody()

	if rrc.errors.Size() > 0 {
		return rrc.errors
	}
	return nil
}

func (rrc *RestResponse) writeHeader() {
	rrc.setHeaderEntries()
	rrc.writeHeaderContent()
}

func (rrc *RestResponse) setHeaderEntries() {
	rrc.setCacheControlHeaderEntry()
	rrc.setContentTypeHeaderEntry()
}

func (rrc *RestResponse) writeHeaderContent() {
	rrc.Writer.WriteHeader(rrc.ResponseCode)
}

func (rrc *RestResponse) setCacheControlHeaderEntry() {
	if rrc.CacheControl != "" {
		rrc.Writer.Header().Set(CacheControlHeaderKey, rrc.CacheControl)
	}
}

func (rrc *RestResponse) setContentTypeHeaderEntry() {
	rrc.Writer.Header().Set(ContentTypeHeaderKey, rrc.ContentType)
}

func (rrc *RestResponse) writeBody() {
	fmt.Fprintf(rrc.Writer, rrc.Content)
}
