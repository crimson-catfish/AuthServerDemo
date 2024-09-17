package handlers

import (
	"TestTask/auth"
	"TestTask/database"
	"TestTask/models"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type loginCredentials struct {
	Guid     string `json:"guid"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	logCred := loginCredentials{}
	err := json.NewDecoder(r.Body).Decode(&logCred)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usr, err := database.GetUserByGuid(logCred.Guid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if usr == (models.User{}) {
		http.Error(w, "No user with such guid is registered", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.HashedPassword), []byte(logCred.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
	}

	expirationTime := time.Now().Add(auth.RefreshTokenExpirationTimeDays * time.Hour * 24)

	err = database.UpdateRefreshToken(usr.Guid, string(hashedRefresh), r.RemoteAddr, expirationTime.Unix())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := registerResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	jResponse, err := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jResponse)
}
