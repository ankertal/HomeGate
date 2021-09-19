package server

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// ServerConfig contains required params
type ServerConfig struct {
	// The http port where the server listens for requests.
	Port string `envconfig:"http_port" default:"80"`

	// database user name
	DBUser string `envconfig:"psql_user" default:"postgres"`

	// database password
	DBPassword string `envconfig:"psql_password" default:""`

	// database name
	DBName string `envconfig:"psql_dbname" default:"postgres"`

	// database host
	DBHost string `envconfig:"psql_host" default:"localhost"`

	// database host
	DBPort string `envconfig:"psql_port" default:"5432"`

	// jwt secret key to sign the tokens
	JWTSecretKey string `envconfig:"jwt_secret_key" default:"homegate;5060102"`

	WebDistro string `envconfig:"web_distro" default:"./frontend/dist/"`
}

// LoadConfig loads  server configuration from the environment
// and validates it.
func LoadConfig() (*ServerConfig, error) {
	conf := ServerConfig{}

	err := envconfig.Process("homegate", &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// Duration type used for json custom marshaling
type Duration struct {
	time.Duration
}

// MarshalJSON duration
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

func (srv *HomeGateServer) setupGates() {
	connection := GetDatabase()
	defer CloseDatabase(connection)

	var gatesDB []Gate
	connection.Table("gates").Find(&gatesDB)

	// we initialize our gates in-memory state
	srv.Lock()
	defer srv.Unlock()

	// setup the gates when we start
	// NOTE: this should probably be revisited to improve perf (user better DB schema etc...)
	for _, gDB := range gatesDB {
		g := &userGate{
			name:       gDB.Name,
			userEmails: make(map[string]bool),
			rcState:    nil,
		}

		// keep the users emails in a map for quick access
		for _, userEmail := range gDB.UserEmails {
			g.userEmails[userEmail] = true
		}

		// keep it in the global gate map
		srv.gates[gDB.Name] = g
	}
}
