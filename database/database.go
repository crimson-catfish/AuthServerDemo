package database

import (
	"TestTask/models"
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (" + models.UserSchema + ");")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")

	DB = db
}

func GetUserByGuid(guid string) (models.User, error) {
	var user models.User

	query := "SELECT * FROM users WHERE guid = $1"
	err := DB.Get(&user, query, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, nil
		}

		return models.User{}, err
	}

	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	query := "SELECT * FROM users WHERE email = $1"
	err := DB.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, nil
		}

		return models.User{}, err
	}

	return user, nil
}

func AddUser(user models.User) error {
	query := "INSERT INTO users (guid, email, hashed_password, last_ip, hashed_refresh_token, refresh_token_expires_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := DB.Exec(query, user.Guid, user.Email, user.HashedPassword, user.LastIP, user.HashedRefreshToken, user.RefreshTokenExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRefreshToken(guid, hashedRefreshToken, ip string, expiresAt int64) error {
	query := `UPDATE users SET hashed_refresh_token=$1, last_ip=$2, refresh_token_expires_at=$3 WHERE guid = $4`

	_, err := DB.Exec(query, hashedRefreshToken, ip, expiresAt, guid)
	if err != nil {
		return err
	}

	return nil
}
