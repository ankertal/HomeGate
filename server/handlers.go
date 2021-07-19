package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var version = "1.0"

type GateEvent struct {
	Deployment *string `json:"deployment,omitempty"`
	User       *string `json:"user,omitempty"`
	Password   *string `json:"password,omitempty"`
}

type CloseEvent struct {
	GateEvent
}

type StatusEvent struct {
	GateEvent
}

type LearnEvent struct {
	GateEvent
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

	//dumpRequest(w, r)

	fmt.Fprintf(w, "HomeGate Server @ %v, version: %v\n", time.Now(), version)
}

func (srv *HomeGateServer) times(w http.ResponseWriter, r *http.Request) {
	srv.Lock()
	defer srv.Unlock()

	//dumpRequest(w, r)

	params := mux.Vars(r)
	deploymentName := params["deployment"]
	dep, ok := srv.deployments[deploymentName]
	if ok {
		fmt.Fprintf(w, "Last Open for Deployment : %v --> %v\n", deploymentName, dep.lastOpen)
		fmt.Fprintf(w, "Last Close for Deployment: %v --> %v\n", deploymentName, dep.lastClose)
		fmt.Fprintf(w, "Last Got Command for RC  : %v --> %v\n", deploymentName, dep.lastGotCommand)
	} else {
		fmt.Fprintf(w, "Could not find a deployment: %v", deploymentName)
	}
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
	//dumpRequest(w, r)

	w.Header().Set("Content-Type", "application/json")
	var openEvent GateEvent
	err := json.NewDecoder(r.Body).Decode(&openEvent)
	if err != nil {
		log.Printf("handler: failed to decode the openEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not openEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	deployment, err := srv.checkGateRequestParams(w, openEvent.Deployment, openEvent.User, openEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	deployment.rcState = Open
	deployment.lastOpen = time.Now()

	fmt.Fprintf(w, "%v's gate requested to OPEN Acknowledged!\n", *deployment.name)
}

func (srv *HomeGateServer) close(w http.ResponseWriter, r *http.Request) {
	//dumpRequest(w, r)

	w.Header().Set("Content-Type", "application/json")
	var closeEvent CloseEvent
	err := json.NewDecoder(r.Body).Decode(&closeEvent)
	if err != nil {
		log.Printf("handler: failed to decode the closeEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not closeEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	deployment, err := srv.checkGateRequestParams(w, closeEvent.Deployment, closeEvent.User, closeEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	deployment.rcState = Close
	deployment.lastClose = time.Now()
	fmt.Fprintf(w, "%v's gate requested to CLOSE Acknowledged!\n", *deployment.name)
}

func (srv *HomeGateServer) rcStatus(w http.ResponseWriter, r *http.Request) {
	//dumpRequest(w, r)

	w.Header().Set("Content-Type", "application/json")
	var statusEvent StatusEvent
	err := json.NewDecoder(r.Body).Decode(&statusEvent)
	if err != nil {
		log.Printf("handler: failed to decode the StatusEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not StatusEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	deployment, err := srv.checkGateRequestParams(w, statusEvent.Deployment, statusEvent.User, statusEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": fmt.Sprintf("%v", deployment.rcState)})

	if deployment.rcState == Open || deployment.rcState == Close || deployment.rcState == Update {
		deployment.lastGotCommand = time.Now()
	}

	deployment.rcState = Unknown
}

func (srv *HomeGateServer) learnOpen(w http.ResponseWriter, r *http.Request) {
	//dumpRequest(w, r)

	w.Header().Set("Content-Type", "application/json")
	var learnEvent LearnEvent
	err := json.NewDecoder(r.Body).Decode(&learnEvent)
	if err != nil {
		log.Printf("handler: failed to decode the learnEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not learnEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	deployment, err := srv.checkGateRequestParams(w, learnEvent.Deployment, learnEvent.User, learnEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	deployment.rcState = LearnOpen

	fmt.Fprintf(w, "%v's gate requested to LEARN Open button -  Acknowledged!\n", *deployment.name)
}

func (srv *HomeGateServer) learnClose(w http.ResponseWriter, r *http.Request) {
	//dumpRequest(w, r)

	w.Header().Set("Content-Type", "application/json")
	var learnEvent LearnEvent
	err := json.NewDecoder(r.Body).Decode(&learnEvent)
	if err != nil {
		log.Printf("handler: failed to decode the learnEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not learnEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	deployment, err := srv.checkGateRequestParams(w, learnEvent.Deployment, learnEvent.User, learnEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	deployment.rcState = LearnClose

	fmt.Fprintf(w, "%v's gate requested to LEARN Close button -  Acknowledged!\n", *deployment.name)
}
