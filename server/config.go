package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
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

// ServerConfig contains required params
type ServerConfig struct {
	// The http port where the server listens for requests.
	Port string `envconfig:"http_port" default:"80"`
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

func (srv *HomeGateServer) setupDeployments() {
	jsonFile, err := os.Open("/Users/yaronweinsberg/work/HomeGate/server/deployments.json")
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
		dep.users = make(map[string]*DeploymentUser)

		for _, user := range configDeployment.Users {
			username := user.Name
			password := user.Password
			dep.users[*username] = &DeploymentUser{Name: username, Password: password}
		}
		srv.deployments[*dep.name] = &dep
	}
}
