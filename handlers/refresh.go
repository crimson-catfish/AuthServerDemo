package handlers

import (
	"TestTask/auth"
	"TestTask/database"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type refreshCredentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	refreshCred := refreshCredentials{}
	err := json.NewDecoder(r.Body).Decode(&refreshCred)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims := auth.AccessTokenClaims{}
	_, err = jwt.ParseWithClaims(refreshCred.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors != jwt.ValidationErrorExpired {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	usr, err := database.GetUserByGuid(claims.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if usr.LastIP != claims.IP {
		//usr.Email.SendMessage("email warning: ip mismatch")
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.HashedRefreshToken), []byte(refreshCred.RefreshToken))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if usr.RefreshTokenExpiresAt < time.Now().Unix() {
		http.Error(w, "Refresh token is expired, please log in with password", http.StatusUnauthorized)
		return
	}

	access, refresh, err := auth.GenerateTokens(usr.Guid, r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expirationTime := time.Now().Add(auth.RefreshTokenExpirationTimeDays * time.Hour * 24)

	err = database.UpdateRefreshToken(usr.Guid, string(hashedRefresh), r.RemoteAddr, expirationTime.Unix())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := refreshResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
