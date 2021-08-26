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

	var users []User
	connection.Preload("Gates").Find(&users)

	// we initialize our gates in-memory state
	srv.Lock()
	defer srv.Unlock()

	gates := make(map[string]*userGate)

	// setup the gates when we start
	// NOTE: this should probably be revisited to improve perf (user better DB schema etc...)
	for _, u := range users {
		g := &userGate{
			name:    u.MyGateName,
			users:   make(map[string]User),
			rcState: make(chan KeyPressed),
		}

		// add this user to its own gate
		g.addUser(u)

		// keep it in the global gate map
		gates[u.MyGateName] = g
	}

	// now since each user can be associated with multiple gates,
	// we update the corresponding gates with the user
	for _, u := range users {
		for _, g := range u.Gates {
			if aGate, ok := gates[g.Name]; ok {
				// this gate is already registered, add this user to (ok to override)
				aGate.addUser(u)
			}
		}
	}

	srv.gates = gates
}
