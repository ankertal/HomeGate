package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
}

type KeyPressed int

const (
	Unknown KeyPressed = iota
	Close
	Open
)

// HomeGateServer represents the webhook server
type HomeGateServer struct {
	config *ServerConfig
	*http.Server
	*mux.Router
	*sync.Mutex
	wg                *sync.WaitGroup
	ShutdownChannel   chan struct{} // shutdown channel
	deploymentRCState map[string]KeyPressed
}

// NewServer creates a new webhook server
func NewServer(config *ServerConfig) *HomeGateServer {
	addr := fmt.Sprintf(":%s", config.Port)
	r := mux.NewRouter()

	srv := &HomeGateServer{
		config:            config,
		Server:            &http.Server{Addr: addr, Handler: r},
		Router:            r,
		Mutex:             &sync.Mutex{},
		ShutdownChannel:   make(chan struct{}),
		wg:                &sync.WaitGroup{},
		deploymentRCState: make(map[string]KeyPressed),
	}

	srv.setupDeployments()

	srv.setupRoutes(r)

	return srv
}

// Shutdown the webhook server
func (srv *HomeGateServer) Shutdown() {
	log.Infof("****** Calling Shutdown ******")
	close(srv.ShutdownChannel)
	srv.Server.Shutdown(context.Background())
}

func (srv *HomeGateServer) setupRoutes(r *mux.Router) {
	srv.Router.HandleFunc("/", srv.home).Methods("GET")
	srv.Router.HandleFunc("/open", srv.open).Methods("POST")
	srv.Router.HandleFunc("/close", srv.close).Methods("POST")
	srv.Router.HandleFunc("/status", srv.rcStatus).Methods("POST")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (srv *HomeGateServer) setupDeployments() {
	for _, deployemnt := range srv.config.Deployments {
		srv.deploymentRCState[deployemnt] = Unknown
	}
}
