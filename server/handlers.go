package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	log "github.com/sirupsen/logrus"
)

var version = "1.0"

type OpenEvent struct {
	Deployment *string `json:"deployment,omitempty"`
	User       *string `json:"user,omitempty"`
	Password   *string `json:"password,omitempty"`
}

type CloseEvent struct {
	OpenEvent
}

type StatusEvent struct {
	OpenEvent
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

	fmt.Fprintf(w, "HomeGate Server @ %v, version: %v", time.Now(), version)
}

func (srv *HomeGateServer) checkGateRequestParams(w http.ResponseWriter, deploymentName, userName, password *string) (*deployment, error) {
	if deploymentName == nil || userName == nil || password == nil {
		http.Error(w, "Bad Request / missing parameter !!!", http.StatusBadRequest)
		return nil, fmt.Errorf("handler: Open called, missing parmeters")
	}

	deployment, ok := srv.deployments[*deploymentName]
	if !ok {
		http.Error(w, "Bad Request / deployment does not exist  !!!", http.StatusBadRequest)
		return nil, fmt.Errorf("handler: unknown deployment: %v", *deploymentName)
	} else {
		user, ok := deployment.users[*userName]
		if !ok {
			http.Error(w, "Bad Request / user does not exist  !!!", http.StatusBadRequest)
			return nil, fmt.Errorf("handler: deployment [ %v ], unknown user: %v", *deploymentName, userName)
		}
		if *user.Password != *password {
			http.Error(w, "Bad Request / user/pass mismatch  !!!", http.StatusBadRequest)
			return nil, fmt.Errorf("handler: deployment [ %v ], user: %v, password mismatch", *deploymentName, userName)
		} else {
			return deployment, nil
		}
	}

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

	deployment, err := srv.checkGateRequestParams(w, openEvent.Deployment, openEvent.User, openEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	deployment.rcState = Open

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

	deployment, err := srv.checkGateRequestParams(w, closeEvent.Deployment, closeEvent.User, closeEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	deployment.rcState = Open

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
		deployment, err := srv.checkGateRequestParams(w, statusEvent.Deployment, statusEvent.User, statusEvent.Password)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{"status": fmt.Sprintf("%v", deployment.rcState)})
		deployment.rcState = Unknown
	}
}
