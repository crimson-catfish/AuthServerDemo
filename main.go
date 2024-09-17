package main

import (
	"TestTask/database"
	"TestTask/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	database.ConnectToDatabase()

	http.HandleFunc("/register", handlers.HandleRegister)
	http.HandleFunc("/login", handlers.HandleLogin)
	http.HandleFunc("/", handlers.HandleHome)

	err := http.ListenAndServe(":"+os.Getenv("LOCALHOST_PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
