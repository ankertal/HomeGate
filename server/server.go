package server

import (
	"fmt"
	"net/http"
	"sync"

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

// userGate represents a user's gate  (identified by the gate uuid name)
type userGate struct {
	name       string
	userEmails map[string]bool
	rcState    chan KeyPressed
}

// HomeGateServer represents the webhook server
type HomeGateServer struct {
	config *ServerConfig
	*http.Server
	*mux.Router
	*sync.Mutex
	ShutdownChannel chan struct{} // shutdown channel
	gates           map[string]*userGate
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
		gates:           make(map[string]*userGate),
	}

	srv.setupRoutes(r)

	InitialMigration()

	srv.setupGates()

	return srv
}

func (srv *HomeGateServer) setupRoutes(r *mux.Router) {
	r.HandleFunc("/command", IsAuthorized(srv.command)).Methods("POST")
	r.HandleFunc("/siri", srv.siri).Methods("POST")
	r.HandleFunc("/stream", srv.stream).Methods("GET")

	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	staticDir := fmt.Sprintf("%v../../static/", srv.config.WebDistro)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	r.HandleFunc("/signup", srv.signUp).Methods("POST")
	r.HandleFunc("/signin", srv.signIn).Methods("POST")
	r.HandleFunc("/user", IsAuthorized(srv.userIndex)).Methods("GET")

	// disable 404 page not found when user refresh screen (SPA known issue)
	r.HandleFunc("/login", redirectHomeFunc).Methods("GET")
	r.HandleFunc("/profile", redirectHomeFunc).Methods("GET")
	r.HandleFunc("/register", redirectHomeFunc).Methods("GET")

	// MUST put this last as order matters
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(srv.config.WebDistro))))
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

var redirectHomeFunc = func(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, fmt.Sprintf("http://%s", req.Host), http.StatusPermanentRedirect)
}
