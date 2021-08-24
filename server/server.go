package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/rs/cors"
)

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
}

type KeyPressed int

const (
	NoOp KeyPressed = iota
	Close
	Open
	Stop
	Update
	LearnOpen
	LearnClose
	LearnStop
	TestOpen
	TestClose
	TestStop
	SetOpen
	SetClose
	SetStop
)

type deployment struct {
	name      *string
	users     map[string]*DeploymentUser
	rcState   chan KeyPressed
	lastOpen  time.Time
	lastClose time.Time
}

// HomeGateServer represents the webhook server
type HomeGateServer struct {
	config *ServerConfig
	*http.Server
	*mux.Router
	*sync.Mutex
	ShutdownChannel chan struct{} // shutdown channel
	deployments     map[string]*deployment
}

// NewServer creates a new webhook server
func NewServer(config *ServerConfig) *HomeGateServer {
	addr := fmt.Sprintf(":%s", config.Port)
	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	srv := &HomeGateServer{
		config: config,
		Server: &http.Server{
			Addr:    addr,
			Handler: c.Handler(r),
		},
		Router:          r,
		Mutex:           &sync.Mutex{},
		ShutdownChannel: make(chan struct{}),
		deployments:     make(map[string]*deployment),
	}

	srv.setupRoutes(r)

	// TODO: this should be populated from a DB
	// TODO: deployment should be also created dynamically via GUI
	srv.setupDeployments()

	return srv
}

func (srv *HomeGateServer) setupRoutes(r *mux.Router) {
	r.HandleFunc("/times/{deployment}", srv.times).Methods("GET")
	r.HandleFunc("/open", srv.open).Methods("POST")
	r.HandleFunc("/close", srv.close).Methods("POST")
	r.HandleFunc("/learn-open", srv.learnOpen).Methods("POST")
	r.HandleFunc("/learn-close", srv.learnClose).Methods("POST")
	r.HandleFunc("/test-open", srv.testOpen).Methods("POST")
	r.HandleFunc("/test-close", srv.testClose).Methods("POST")
	r.HandleFunc("/set-open", srv.setOpen).Methods("POST")
	r.HandleFunc("/set-close", srv.setClose).Methods("POST")
	r.HandleFunc("/stream", srv.stream).Methods("GET")

	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/signup", srv.signUp).Methods("POST")
	r.HandleFunc("/signin", srv.signIn).Methods("POST")
	r.HandleFunc("/admin", IsAuthorized(srv.adminIndex)).Methods("GET")
	r.HandleFunc("/user", IsAuthorized(srv.userIndex)).Methods("GET")

	// MUST put this last as order matters
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend/dist/"))))

	InitialMigration()
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})
