package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// check whether user is authorized or not
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			var err Error
			err = SetError(err, "No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte(os.Getenv("HOMEGATE_JWT_SECRET_KEY"))

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error in parsing token")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err Error
			err = SetError(err, "Your Token has been expired.")
			err.sendToClient(w, http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return

			} else if claims["role"] == "user" {
				userEmail, ok := claims["email"]
				if ok && userEmail != "" {
					r.Header.Set("Role", "user")
					r.Header.Set("Email", userEmail.(string))
					handler.ServeHTTP(w, r)
				} else {
					var err Error
					err = SetError(err, "Your Token is bogus, please login again")
					err.sendToClient(w, http.StatusUnauthorized)
					return
				}
				return

			}
		}
		var reserr Error
		reserr = SetError(reserr, "Not Authorized.")
		reserr.sendToClient(w, http.StatusUnauthorized)
	}
}
