package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	log "github.com/sirupsen/logrus"
)

type OpenEvent struct {
	Deployment *string `json:"deployment,omitempty"`
	User       *string `json:"user,omitempty"`
	Password   *string `json:"password,omitempty"`
}

type CloseEvent struct {
	OpenEvent
}

type StatusEvent struct {
	Deployment *string `json:"deployment,omitempty"`
	Password   *string `json:"password,omitempty"`
}

func dumpRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler: inside hooksHandler")
	var err error

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	log.Printf(string(requestDump))

}

func (srv *HomeGateServer) home(w http.ResponseWriter, r *http.Request) {
	srv.Lock()
	defer srv.Unlock()

	dumpRequest(w, r)

	fmt.Fprintf(w, "Hello from HomeGate Server @ %v", time.Now())
}

func (srv *HomeGateServer) open(w http.ResponseWriter, r *http.Request) {
	srv.Lock()
	defer srv.Unlock()

	dumpRequest(w, r)

	w.Header().Set("Content-Type", "application/json")
	var openEvent OpenEvent
	err := json.NewDecoder(r.Body).Decode(&openEvent)
	if err != nil {
		log.Printf("handler: failed to decode the openEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not openEvent !!!", http.StatusBadRequest)
		return
	}
	srv.deploymentRCState[*openEvent.Deployment] = Open

}

func (srv *HomeGateServer) close(w http.ResponseWriter, r *http.Request) {
	srv.Lock()
	defer srv.Unlock()

	dumpRequest(w, r)

	w.Header().Set("Content-Type", "application/json")
	var closeEvent CloseEvent
	err := json.NewDecoder(r.Body).Decode(&closeEvent)
	if err != nil {
		log.Printf("handler: failed to decode the closeEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not closeEvent !!!", http.StatusBadRequest)
		return
	}
	srv.deploymentRCState[*closeEvent.Deployment] = Open

}

func (srv *HomeGateServer) rcStatus(w http.ResponseWriter, r *http.Request) {
	srv.Lock()
	defer srv.Unlock()

	dumpRequest(w, r)

	w.Header().Set("Content-Type", "application/json")
	var statusEvent StatusEvent
	err := json.NewDecoder(r.Body).Decode(&statusEvent)
	if err != nil {
		log.Printf("handler: failed to decode the StatusEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not StatusEvent !!!", http.StatusBadRequest)
		return
	}
	if true {
		respondWithJSON(w, http.StatusOK, map[string]string{"status": fmt.Sprintf("%v", srv.deploymentRCState[*statusEvent.Deployment])})
		srv.deploymentRCState[*statusEvent.Deployment] = Unknown
	}
}
