package handlers

import (
	"TestTask/auth"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
)

type homeCredential struct {
	AccessToken string `json:"access_token"`
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	homeCred := homeCredential{}
	err := json.NewDecoder(r.Body).Decode(&homeCred)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(homeCred.AccessToken)

	claims := auth.AccessTokenClaims{}
	_, err = jwt.ParseWithClaims(homeCred.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if claims.IP != r.RemoteAddr {
		http.Error(w, "IP address mismatch, please refresh token", http.StatusForbidden)
		return
	}

	w.Write([]byte("Access confirmed\nWelcome to the home page!\n"))
}
