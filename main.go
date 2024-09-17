package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	ConnectToDatabase()

	http.HandleFunc("/register", HandleRegister)

	err := http.ListenAndServe(":"+os.Getenv("LOCALHOST_PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
