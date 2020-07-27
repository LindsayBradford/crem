// Copyright (c) 2018 Australian Rivers Institute.

package admin

import (
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
)

type ServiceStatus struct {
	ServiceName string
	Version     string
	Status      string
	Time        string
}

const muxType = "ADMIN"

type Mux struct {
	rest.MuxImpl
	Status ServiceStatus

	doneChannel chan bool
}

func (m *Mux) Initialise() *Mux {
	m.MuxImpl.Initialise().WithType(muxType)

	m.doneChannel = make(chan bool)
	m.AddHandler("/status", m.StatusHandler)
	m.AddHandler("/shutdown", m.shutdownHandler)

	return m
}

func (m *Mux) AddHandler(addressPattern string, handler rest.HandlerFunc) {
	m.HandlerMap.AddHandler(addressPattern, handler)
}

func (m *Mux) WithType(muxType string) *Mux {
	m.MuxImpl.WithType(muxType)
	return m
}

func (m *Mux) SetStatus(statusMessage string) {
	m.Logger().Info("Changing server status to [" + statusMessage + "]")
	m.Status.Status = statusMessage
	m.UpdateStatusTime()
}

func (m *Mux) Start(address string) {
	m.MuxImpl.Start(address)
}

func (m *Mux) WaitForShutdownSignal() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, os.Kill)
		<-sigint

		m.Logger().Warn("Received Operating System Interrupt/Kill signal -- triggering graceful shutdown")

		close(m.doneChannel)
	}()

	<-m.doneChannel
}

func (m *Mux) StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		m.MethodNotAllowedError(w, r)
		return
	}

	m.UpdateStatusTime()
	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(m.Status)

	m.Logger().Info("Responding with status [" + m.Status.Status + "]")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "status handler")
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		m.MethodNotAllowedError(w, r)
		return
	}

	m.Status.Status = "SHUTTING_DOWN"
	m.UpdateStatusTime()

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(m.Status)

	m.Logger().Debug("Responding with status [" + m.Status.Status + "]")
	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "shutdown handler")
		m.Logger().Error(wrappingError)
	}

	m.doneChannel <- true
}

func (m *Mux) UpdateStatusTime() {
	m.Status.Time = rest.FormattedTimestamp()
}

func (m *Mux) Shutdown() {
	m.Server().Shutdown(context.Background())
	m.SetStatus("DEAD")
}
