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

type registerCredentials struct {
	Guid     string `json:"guid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	regCred := registerCredentials{}
	err := json.NewDecoder(r.Body).Decode(&regCred)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usr, err := database.GetUserByGuid(regCred.Guid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if usr != (models.User{}) {
		http.Error(w, "User with such guid already exists", http.StatusConflict)
		return
	}

	usr, err = database.GetUserByEmail(regCred.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if usr != (models.User{}) {
		http.Error(w, "User with such email already exists", http.StatusConflict)
		return
	}

	access, refresh, err := auth.GenerateTokens(usr.Guid, r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(regCred.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	expirationTime := time.Now().Add(auth.RefreshTokenExpirationTimeDays * time.Hour * 24)

	usr = models.User{
		Guid:                  regCred.Guid,
		Email:                 regCred.Email,
		HashedPassword:        string(hashedPass),
		LastIP:                r.RemoteAddr,
		HashedRefreshToken:    string(hashedRefresh),
		RefreshTokenExpiresAt: expirationTime.Unix(),
	}

	err = database.AddUser(usr)
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
