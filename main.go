package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	ConnectToDatabase()

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!"))
	})

	err := http.ListenAndServe(":"+os.Getenv("LOCALHOST_PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
