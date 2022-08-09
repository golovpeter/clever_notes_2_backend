package main

import (
	"clever_notes_2/internal/handlers/signup"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/signup", signup.SingUp)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
