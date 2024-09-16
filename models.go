package main

type User struct {
	Guid               string `json:"guid" db:"guid"`
	Email              string `json:"email" db:"email"`
	HashedPassword     string `json:"hashed_password" db:"hashed_password"`
	LastIP             string `json:"last_ip" db:"last_ip"`
	HashedRefreshToken string `json:"hashed_refresh_token" db:"hashed_refresh_token"`
}

var UserSchema = `
	guid VARCHAR(36) PRIMARY KEY,
	email TEXT NOT NULL,
	hashed_password TEXT NOT NULL,
	last_ip VARCHAR(45) NOT NULL,
	hashed_refresh_token TEXT NOT NULL`
