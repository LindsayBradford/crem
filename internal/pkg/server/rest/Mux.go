// Copyright (c) 2018 Australian Rivers Institute.

package rest

import (
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type Mux interface {
	Start(portAddress string)
	Shutdown()

	SetLogger(handler logging.Logger)
	AddHandler(address string, handler HandlerFunc)

	SetCacheMaxAge(maxAge uint64)
	CacheMaxAge() uint64
}

const DefaultCacheMaxAgeInSeconds = 10

type MuxImpl struct {
	http.ServeMux
	muxType              string
	server               http.Server
	cacheMaxAgeInSeconds uint64

	HandlerMap HandlerFunctionMap
	logger     logging.Logger
}

type HandlerFunctionMap map[string]HandlerFunc

func (mi *MuxImpl) Initialise() *MuxImpl {
	mi.HandlerMap = make(HandlerFunctionMap)
	mi.SetCacheMaxAge(DefaultCacheMaxAgeInSeconds)
	return mi
}

func (mi *MuxImpl) WithType(muxType string) *MuxImpl {
	mi.muxType = muxType
	return mi
}

func (mi *MuxImpl) WithCacheMaxAge(maxAgeInSeconds uint64) *MuxImpl {
	if maxAgeInSeconds != 0 {
		mi.cacheMaxAgeInSeconds = maxAgeInSeconds
	}
	return mi
}

func (mi *MuxImpl) SetLogger(logger logging.Logger) {
	mi.logger = logger
}

func (mi *MuxImpl) SetCacheMaxAge(maxAgeInSeconds uint64) {
	mi.cacheMaxAgeInSeconds = maxAgeInSeconds
}

func (mi *MuxImpl) CacheMaxAge() uint64 {
	return mi.cacheMaxAgeInSeconds
}

func (mi *MuxImpl) Logger() logging.Logger {
	return mi.logger
}

func (mi *MuxImpl) Server() *http.Server {
	return &mi.server
}

func (mi *MuxImpl) AddHandler(address string, handler HandlerFunc) {
	mi.HandlerMap[address] = handler
}

func (mi *MuxImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mi.logRequestReceipt(r)
	if handlerFunction, handlerFound := mi.handlerFor(r); handlerFound {
		handlerFunction(w, r)
	} else {
		mi.NotFoundError(w, r)
	}
}

func (mi *MuxImpl) logRequestReceipt(r *http.Request) {
	mi.logger.Info(
		"[" + mi.muxType + "] Received request Method [" + r.Method +
			"] for request [" + r.URL.Path + "] from [" + r.RemoteAddr + "].")
}

func (mi *MuxImpl) handlerFor(r *http.Request) (handlerFunction HandlerFunc, found bool) {
	handlerFunction, found = mi.HandlerMap[r.URL.String()]
	return
}

func (mi *MuxImpl) NotFoundError(w http.ResponseWriter, r *http.Request) {
	mi.RespondWithError(http.StatusNotFound, "HTTP Resource not found", w, r)
}

func (mi *MuxImpl) MethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	mi.RespondWithError(http.StatusMethodNotAllowed, "HTTP Method not allowed", w, r)
}

func (mi *MuxImpl) InternalServerError(w http.ResponseWriter, r *http.Request, errorDetail error) {
	finalErrorString := "Internal Server Error"
	if errorDetail != nil {
		finalErrorString = fmt.Sprintf("%s: %v", finalErrorString, errorDetail)
	}
	mi.RespondWithError(http.StatusInternalServerError, finalErrorString, w, r)
}

func (mi *MuxImpl) ServiceUnavailableError(w http.ResponseWriter, r *http.Request, errorDetail error) {
	finalErrorString := "Service Unavailable Error"
	if errorDetail != nil {
		finalErrorString = fmt.Sprintf("%s: %v", finalErrorString, errorDetail)
	}
	mi.RespondWithError(http.StatusServiceUnavailable, finalErrorString, w, r)
}

func (mi *MuxImpl) RespondWithError(responseCode int, responseMsg string, w http.ResponseWriter, r *http.Request) {
	mi.logResponseError(r, responseMsg)

	restResponse := new(Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(responseCode).
		WithJsonContent(
			ErrorResponse{ErrorMessage: responseMsg, Time: FormattedTimestamp()},
		)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "responding with error")
		mi.logger.Error(wrappingError)
	}
}

func (mi *MuxImpl) logResponseError(r *http.Request, responseMsg string) {
	mi.logger.Warn(
		"Request Method [" + r.Method + "] for request [" + r.URL.Path + "] from [" + r.RemoteAddr +
			"]. Responding with [" + responseMsg + "] error.")
}

func (mi *MuxImpl) Start(address string) {
	mi.logger.Debug("Starting [" + mi.muxType + "] server on address [" + address + "]")

	mi.server = http.Server{Addr: address, Handler: mi}

	if err := mi.server.ListenAndServe(); err != http.ErrServerClosed {
		wrappedErr := errors.Wrap(err, "ListenAndServe")
		mi.logger.Error(wrappedErr)
	}
}

func (mi *MuxImpl) Shutdown() {
	mi.server.Shutdown(context.Background())
}
