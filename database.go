package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sqlx.DB

func ConnectToDatabase() {
	connectStr := fmt.Sprintf(
		"host=DB user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Connect("postgres", connectStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (" + UserSchema + ");")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")

	DB = db
}

func GetUserByGuid(guid string) (User, error) {
	var user User

	query := "SELECT * FROM users WHERE guid = $1"
	err := DB.Get(&user, query, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}

		return User{}, err
	}

	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User

	query := "SELECT * FROM users WHERE guid = $1"
	err := DB.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}

		return User{}, err
	}

	return user, nil
}

func AddUser(user User) error {
	query := "INSERT INTO users (guid, email, hashed_password, last_ip, hashed_refresh_token) VALUES ($1, $2, $3, $4, $5)"

	_, err := DB.Exec(query, user.Guid, user.Email, user.HashedPassword, user.LastIP, user.HashedRefreshToken)
	if err != nil {
		return err
	}

	return nil
}
