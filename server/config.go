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
	Port        string   `envconfig:"http_port" default:"80"`
	Deployments []string `envconfig:"deployments" default:"Taron,Tal,Gilad,Doron"`
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
