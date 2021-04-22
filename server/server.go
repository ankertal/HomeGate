package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	Stop
	Update
)

type DeploymentUser struct {
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
}

type DeploymentConfig struct {
	Name  *string          `json:"name,omitempty"`
	Users []DeploymentUser `json:"users"`
}

type DeploymentsConfig struct {
	Deployments []DeploymentConfig `json:"deployments"`
}

type deployment struct {
	name    *string
	users   map[string]*DeploymentUser
	rcState KeyPressed
}

// HomeGateServer represents the webhook server
type HomeGateServer struct {
	config *ServerConfig
	*http.Server
	*mux.Router
	*sync.Mutex
	wg              *sync.WaitGroup
	ShutdownChannel chan struct{} // shutdown channel
	deployments     map[string]*deployment
}

// NewServer creates a new webhook server
func NewServer(config *ServerConfig) *HomeGateServer {
	addr := fmt.Sprintf(":%s", config.Port)
	r := mux.NewRouter()

	srv := &HomeGateServer{
		config:          config,
		Server:          &http.Server{Addr: addr, Handler: r},
		Router:          r,
		Mutex:           &sync.Mutex{},
		ShutdownChannel: make(chan struct{}),
		wg:              &sync.WaitGroup{},
		deployments:     make(map[string]*deployment),
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
	jsonFile, err := os.Open("/home/ankertal/Work/HomeGate/server/deployments.json")
	if err != nil {
		panic("Could not find a deployments file")
	}

	fmt.Println("Successfully Opened deployments.json")
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configDeployments DeploymentsConfig
	json.Unmarshal(byteValue, &configDeployments)

	// we initialize our deployments 'button' states
	for _, configDeployment := range configDeployments.Deployments {
		var dep deployment
		dep.name = configDeployment.Name
		dep.rcState = Update
		dep.users = make(map[string]*DeploymentUser)

		for _, user := range configDeployment.Users {
			username := user.Name
			password := user.Password
			dep.users[*username] = &DeploymentUser{Name: username, Password: password}
		}
		srv.deployments[*dep.name] = &dep
	}
}
