package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

	srv := &HomeGateServer{
		config: config,
		Server: &http.Server{
			Addr: addr,
			Handler: handlers.CORS(
				handlers.AllowedHeaders(
					[]string{
						"X-Requested-With",
						"Access-Control-Allow-Origin",
						"Content-Type", "Authorization",
					}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{"*"}))(r),
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
	srv.Router.HandleFunc("/times/{deployment}", srv.times).Methods("GET")
	srv.Router.HandleFunc("/open", srv.open).Methods("POST")
	srv.Router.HandleFunc("/close", srv.close).Methods("POST")
	srv.Router.HandleFunc("/learn-open", srv.learnOpen).Methods("POST")
	srv.Router.HandleFunc("/learn-close", srv.learnClose).Methods("POST")
	srv.Router.HandleFunc("/test-open", srv.testOpen).Methods("POST")
	srv.Router.HandleFunc("/test-close", srv.testClose).Methods("POST")
	srv.Router.HandleFunc("/set-open", srv.setOpen).Methods("POST")
	srv.Router.HandleFunc("/set-close", srv.setClose).Methods("POST")
	srv.Router.HandleFunc("/stream", srv.stream).Methods("GET")

	// On the default page we will simply serve our static index page.
	r.Handle("/", http.FileServer(http.Dir("./views/")))

	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/signup", srv.signUp).Methods("POST")
	r.HandleFunc("/signin", srv.signIn).Methods("POST")
	r.HandleFunc("/admin", IsAuthorized(srv.adminIndex)).Methods("GET")
	r.HandleFunc("/user", IsAuthorized(srv.userIndex)).Methods("GET")
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})

	InitialMigration()
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})
