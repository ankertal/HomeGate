package server

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Name     string `json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Gates    []Gate `json:"gates" gorm:"many2many:user_gates;"`
}

type Gate struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
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
	Gates       []string `json:"gates"`
}

type RegisterResponse struct {
	Name    string `json:"username"`
	Email   string `gorm:"unique" json:"email"`
	Message string `json:"message"`
}
