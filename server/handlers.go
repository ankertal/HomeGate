package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var version = "1.0"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Resolve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

type TestEvent struct {
	GateEvent
}

type SetEvent struct {
	GateEvent
}

func (srv *HomeGateServer) times(w http.ResponseWriter, r *http.Request) {
	srv.Lock()
	defer srv.Unlock()

	params := mux.Vars(r)
	deploymentName := params["deployment"]
	dep, ok := srv.deployments[deploymentName]
	if ok {
		fmt.Fprintf(w, "Last Open for Deployment : %v --> %v\n", deploymentName, dep.lastOpen)
		fmt.Fprintf(w, "Last Close for Deployment: %v --> %v\n", deploymentName, dep.lastClose)
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

	if deployment.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case deployment.rcState <- Open:
		deployment.lastOpen = time.Now()
		fmt.Fprintf(w, "%v's gate requested to OPEN Acknowledged!\n", *deployment.name)
	default:
		log.Printf("client does not read events...")
		close(deployment.rcState)
	}

}

func (srv *HomeGateServer) close(w http.ResponseWriter, r *http.Request) {

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

	if deployment.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case deployment.rcState <- Close:
		deployment.lastClose = time.Now()
		fmt.Fprintf(w, "%v's gate requested to CLOSE Acknowledged!\n", *deployment.name)
	default:
		log.Printf("client does not read events...")
		close(deployment.rcState)
	}

}

func (srv *HomeGateServer) learnOpen(w http.ResponseWriter, r *http.Request) {

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

	if deployment.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case deployment.rcState <- LearnOpen:
		fmt.Fprintf(w, "%v's gate requested to LearnOpen Acknowledged!\n", *deployment.name)
	default:
		log.Printf("client does not read events...")
		close(deployment.rcState)
	}

	fmt.Fprintf(w, "%v's gate requested to LEARN Open button -  Acknowledged!\n", *deployment.name)
}

func (srv *HomeGateServer) learnClose(w http.ResponseWriter, r *http.Request) {

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

	if deployment.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case deployment.rcState <- LearnClose:
		fmt.Fprintf(w, "%v's gate requested to LearnClose Acknowledged!\n", *deployment.name)
	default:
		log.Printf("client does not read events...")
		close(deployment.rcState)
	}

	fmt.Fprintf(w, "%v's gate requested to LEARN Close button -  Acknowledged!\n", *deployment.name)
}

func (srv *HomeGateServer) testOpen(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var testEvent TestEvent
	err := json.NewDecoder(r.Body).Decode(&testEvent)
	if err != nil {
		log.Printf("handler: failed to decode the testEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not testEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	deployment, err := srv.checkGateRequestParams(w, testEvent.Deployment, testEvent.User, testEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if deployment.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case deployment.rcState <- TestOpen:
		fmt.Fprintf(w, "%v's gate requested to TestOpen Acknowledged!\n", *deployment.name)
	default:
		log.Printf("client does not read events...")
		close(deployment.rcState)
	}
}

func (srv *HomeGateServer) testClose(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var testEvent TestEvent
	err := json.NewDecoder(r.Body).Decode(&testEvent)
	if err != nil {
		log.Printf("handler: failed to decode the testEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not testEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	deployment, err := srv.checkGateRequestParams(w, testEvent.Deployment, testEvent.User, testEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if deployment.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case deployment.rcState <- TestClose:
		fmt.Fprintf(w, "%v's gate requested to TestClose Acknowledged!\n", *deployment.name)
	default:
		log.Printf("client does not read events...")
		close(deployment.rcState)
	}
}

func (srv *HomeGateServer) setOpen(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var setEvent SetEvent
	err := json.NewDecoder(r.Body).Decode(&setEvent)
	if err != nil {
		log.Printf("handler: failed to decode the setEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not testEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	deployment, err := srv.checkGateRequestParams(w, setEvent.Deployment, setEvent.User, setEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if deployment.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case deployment.rcState <- SetOpen:
		fmt.Fprintf(w, "%v's gate requested to SetOpen Acknowledged!\n", *deployment.name)
	default:
		log.Printf("client does not read events...")
		close(deployment.rcState)
	}
}

func (srv *HomeGateServer) setClose(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var setEvent SetEvent
	err := json.NewDecoder(r.Body).Decode(&setEvent)
	if err != nil {
		log.Printf("handler: failed to decode the setEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not testEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	deployment, err := srv.checkGateRequestParams(w, setEvent.Deployment, setEvent.User, setEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if deployment.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case deployment.rcState <- SetClose:
		fmt.Fprintf(w, "%v's gate requested to SetClose Acknowledged!\n", *deployment.name)
	default:
		log.Printf("client does not read events...")
		close(deployment.rcState)
	}
}

func (srv *HomeGateServer) stream(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// validate user
	mt, message, err := c.ReadMessage()
	if err != nil {
		http.Error(w, "Bad Request / data is not StatusEvent !!!", http.StatusBadRequest)
		return
	}
	var statusEvent StatusEvent
	err = json.Unmarshal([]byte(message), &statusEvent)
	if err != nil {
		http.Error(w, "Bad Request / data is not StatusEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	deployment, err := srv.checkGateRequestParams(w, statusEvent.Deployment, statusEvent.User, statusEvent.Password)
	if err != nil {
		http.Error(w, "Bad Request / invalid deployment", http.StatusBadRequest)
		return
	}

	// create a new channel for the deployment
	if deployment.rcState != nil {
		close(deployment.rcState)
	}

	deployment.rcState = make(chan KeyPressed, 1)
	srv.Unlock()

	for rcEvent := range deployment.rcState {
		err = c.WriteMessage(mt, []byte(fmt.Sprintf("%v", rcEvent)))
		if err != nil {
			break
		}
	}
}
