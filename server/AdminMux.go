// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/net/context"
)

type ServiceStatus struct {
	ServiceName string
	Version     string
	Status      string
	Time        string
}

const AdminMuxType = "ADMIN"

type AdminMux struct {
	BaseMux
	Status ServiceStatus

	doneChannel chan bool
}

func (am *AdminMux) Initialise() *AdminMux {
	am.BaseMux.Initialise().WithType(AdminMuxType)

	am.doneChannel = make(chan bool)
	am.handlerMap["/status"] = am.statusHandler
	am.handlerMap["/shutdown"] = am.shutdownHandler

	return am
}

func (am *AdminMux) WithType(muxType string) *AdminMux {
	am.BaseMux.WithType(muxType)
	return am
}

func (am *AdminMux) setStatus(statusMessage string) {
	am.logger.Info("Changed server Status to [" + statusMessage + "]")
	am.Status.Status = statusMessage
	am.UpdateStatusTime()
}

func (am *AdminMux) Start(address string) {
	am.setStatus("RUNNING")
	am.BaseMux.Start(address)
}

func (am *AdminMux) WaitForShutdownSignal() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, os.Kill)
		<-sigint

		am.logger.Warn("Received Operating System Interrupt/Kill signal -- triggering graceful shutdown")

		close(am.doneChannel)
	}()

	<-am.doneChannel
}

func (am *AdminMux) statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		am.MethodNotAllowedError(w, r)
		return
	}

	am.logger.Debug("Responding with status [" + am.Status.Status + "]")
	am.UpdateStatusTime()

	setResponseContentType(w, JsonMimeType)
	am.SetResponseCacheMaxAge(w)

	statusJson, encodeError := json.MarshalIndent(am.Status, "", "  ")
	if encodeError != nil {
		am.logger.Error(encodeError)
	}

	fmt.Fprintf(w, string(statusJson))
}

func (am *AdminMux) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		am.MethodNotAllowedError(w, r)
		return
	}

	am.Status.Status = "SHUTTING_DOWN"
	am.logger.Debug("Responding with status [" + am.Status.Status + "]")
	am.UpdateStatusTime()

	setResponseContentType(w, JsonMimeType)
	am.SetResponseCacheMaxAge(w)

	statusJson, encodeError := json.MarshalIndent(am.Status, "", "  ")
	if encodeError != nil {
		am.logger.Error(encodeError)
	}

	bufferedWriter := bufio.NewWriter(w)

	fmt.Fprintf(bufferedWriter, string(statusJson))
	bufferedWriter.Flush()

	am.doneChannel <- true
}

func (am *AdminMux) UpdateStatusTime() {
	am.Status.Time = FormattedTimestamp()
}

func (am *AdminMux) Shutdown() {
	am.server.Shutdown(context.Background())
	am.setStatus("DEAD")
}
