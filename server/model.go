package server

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

type User struct {
	gorm.Model
	Name       string         `json:"username"`
	Email      string         `gorm:"unique" json:"email"`
	Password   string         `json:"password"`
	Role       string         `json:"role"`
	MyGateName string         `json:"my_gate"`
	Gates      pq.StringArray `json:"gates" gorm:"type:text[]"`
}

type Gate struct {
	gorm.Model
	Name       string         `json:"name" gorm:"unique"`
	UserEmails pq.StringArray `json:"user_emails" gorm:"type:text[]"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

type LoginResponse struct {
	ID          uint     `json:"id"`
	Message     string   `json:"message"`
	UserName    string   `json:"username"`
	Email       string   `json:"email"`
	AccessToken string   `json:"accessToken"`
	MyGateName  string   `json:"my_gate"`
	Gates       []string `json:"gates"`
}

type RegisterResponse struct {
	Name    string `json:"username"`
	Email   string `gorm:"unique" json:"email"`
	Message string `json:"message"`
}

type GateEvent struct {
	GateName     *string `json:"gate_name,omitempty"`
	IsOpen       *bool   `json:"is_open,omitempty"`
	IsClose      *bool   `json:"is_close,omitempty"`
	IsLearnOpen  *bool   `json:"is_learn_open,omitempty"`
	IsLearnClose *bool   `json:"is_learn_close,omitempty"`
	IsTestOpen   *bool   `json:"is_test_open,omitempty"`
	IsTestClose  *bool   `json:"is_test_close,omitempty"`
	IsSetOpen    *bool   `json:"is_set_open,omitempty"`
	IsSetClose   *bool   `json:"is_set_close,omitempty"`
}

// DeviceStreamRequest is the request send from the gate controller (pi) to the stream endpoint
type DeviceStreamRequest struct {
	GateID string `json:"gate_id,omitempty"`
}

// SiriCommand is used for SR command via the phone, it uses simple gate ID (secret) authentication.
type SiriCommand struct {
	GateID             string `json:"gate_id,omitempty"`
	OpenOrCloseCommand string `json:"gate_command,omitempty"`
}
