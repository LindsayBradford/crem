// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crm/logging/handlers"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type RestMux interface {
	Start(portAddress string)
	Shutdown()

	SetLogger(handler handlers.LogHandler)
	AddHandler(address string, handler http.HandlerFunc)

	SetCacheMaxAge(maxAge uint64)
	CacheMaxAge() uint64
}

const DefaultCacheMaxAgeInSeconds = 10

type BaseMux struct {
	http.ServeMux
	muxType              string
	server               http.Server
	cacheMaxAgeInSeconds uint64

	handlerMap HandlerFunctionMap
	logger     handlers.LogHandler
}

type HandlerFunctionMap map[string]http.HandlerFunc

func (bm *BaseMux) Initialise() *BaseMux {
	bm.handlerMap = make(HandlerFunctionMap)
	bm.SetCacheMaxAge(DefaultCacheMaxAgeInSeconds)
	return bm
}

func (bm *BaseMux) WithType(muxType string) *BaseMux {
	bm.muxType = muxType
	return bm
}

func (bm *BaseMux) WithCacheMaxAge(maxAgeInSeconds uint64) *BaseMux {
	if maxAgeInSeconds != 0 {
		bm.cacheMaxAgeInSeconds = maxAgeInSeconds
	}
	return bm
}

func (bm *BaseMux) SetLogger(logger handlers.LogHandler) {
	bm.logger = logger
}

func (bm *BaseMux) SetCacheMaxAge(maxAgeInSeconds uint64) {
	bm.cacheMaxAgeInSeconds = maxAgeInSeconds
}

func (bm *BaseMux) CacheMaxAge() uint64 {
	return bm.cacheMaxAgeInSeconds
}

func (bm *BaseMux) Logger() handlers.LogHandler {
	return bm.logger
}

func (bm *BaseMux) AddHandler(address string, handler http.HandlerFunc) {
	bm.handlerMap[address] = handler
}

func (bm *BaseMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bm.logRequestReceipt(r)
	if handlerFunction, handlerFound := bm.handlerFor(r); handlerFound {
		handlerFunction(w, r)
	} else {
		bm.NotFoundError(w, r)
	}
}

func (bm *BaseMux) logRequestReceipt(r *http.Request) {
	bm.logger.Info(
		"[" + bm.muxType + "] Received request Method [" + r.Method +
			"] for request [" + r.URL.Path + "] from [" + r.RemoteAddr + "].")
}

func (bm *BaseMux) handlerFor(r *http.Request) (handlerFunction http.HandlerFunc, found bool) {
	handlerFunction, found = bm.handlerMap[r.URL.String()]
	return
}

func (bm *BaseMux) NotFoundError(w http.ResponseWriter, r *http.Request) {
	bm.RespondWithError(http.StatusNotFound, "HTTP Resource not found", w, r)
}

func (bm *BaseMux) MethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	bm.RespondWithError(http.StatusMethodNotAllowed, "HTTP Method not allowed", w, r)
}

func (bm *BaseMux) InternalServerError(w http.ResponseWriter, r *http.Request, errorDetail error) {
	finalErrorString := "Internal Server Error"
	if errorDetail != nil {
		finalErrorString = fmt.Sprintf("%s: %v", finalErrorString, errorDetail)
	}
	bm.RespondWithError(http.StatusInternalServerError, finalErrorString, w, r)
}

func (bm *BaseMux) RespondWithError(responseCode int, responseMsg string, w http.ResponseWriter, r *http.Request) {
	bm.logResponseError(r, responseMsg)

	restResponse := new(RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(responseCode).
		WithJsonContent(
			ErrorResponse{ErrorMessage: responseMsg, Time: FormattedTimestamp()},
		)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "responding with error")
		bm.logger.Error(wrappingError)
	}
}

func setResponseContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set(ContentTypeHeaderKey, contentType)
}

func (bm *BaseMux) SetResponseCacheMaxAge(w http.ResponseWriter) {
	maxAgeAsString := fmt.Sprintf("max-age=%d", bm.cacheMaxAgeInSeconds)
	w.Header().Set(CacheControlHeaderKey, maxAgeAsString)
}

func (bm *BaseMux) logResponseError(r *http.Request, responseMsg string) {
	bm.logger.Warn(
		"Request Method [" + r.Method + "] for request [" + r.URL.Path + "] from [" + r.RemoteAddr +
			"]. Responding with [" + responseMsg + "] error.")
}

func (bm *BaseMux) Start(address string) {
	bm.logger.Debug("Starting [" + bm.muxType + "] server on address [" + address + "]")

	bm.server = http.Server{Addr: address, Handler: bm}

	if err := bm.server.ListenAndServe(); err != http.ErrServerClosed {
		wrappedErr := errors.Wrap(err, "ListenAndServe")
		bm.logger.Error(wrappedErr)
	}
}

func (bm *BaseMux) Shutdown() {
	bm.server.Shutdown(context.Background())
}
