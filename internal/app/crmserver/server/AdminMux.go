// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
)

type AdminMux struct {
	RestMux
	Status Status

	doneChannel chan bool
}

func (am *AdminMux) Initialise() *AdminMux {
	am.RestMux.Initialise()

	am.doneChannel = make(chan bool)
	am.handlerMap["/status"] = am.statusHandler
	am.handlerMap["/shutdown"] = am.shutdownHandler

	return am
}

func (am *AdminMux) WithType(muxType string) *AdminMux {
	am.RestMux.WithType(muxType)
	return am
}

func (am *AdminMux) setStatus(statusMessage string) {
	am.Logger.Debug("Changed server Status to [" + statusMessage + "]")
	am.Status.Message = statusMessage
	am.UpdateStatusTime()
}

func (am *AdminMux) WaitOnShutdown() {
	<-am.doneChannel
}

func (am *AdminMux) statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		am.ServeMethodNotAllowedError(w, r)
		return
	}

	am.Logger.Debug("Responding with status [" + am.Status.Message + "]")
	am.UpdateStatusTime()

	statusJson, encodeError := json.MarshalIndent(am.Status, "", "  ")
	if encodeError != nil {
		am.Logger.Error(encodeError)
	}

	fmt.Fprintf(w, string(statusJson))
}

func (am *AdminMux) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		am.ServeMethodNotAllowedError(w, r)
		return
	}

	am.Status.Message = "SHUTTING_DOWN"
	am.Logger.Debug("Responding with status [" + am.Status.Message + "]")
	am.UpdateStatusTime()

	statusJson, encodeError := json.MarshalIndent(am.Status, "", "  ")
	if encodeError != nil {
		am.Logger.Error(encodeError)
	}

	bufferedWriter := bufio.NewWriter(w)

	fmt.Fprintf(bufferedWriter, string(statusJson))
	bufferedWriter.Flush()

	am.doneChannel <- true
}

func (am *AdminMux) UpdateStatusTime() {
	am.Status.Time = FormattedTimestamp()
}
