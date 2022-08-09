package main

import (
	"fmt"
	"github.com/golovpeter/clever_notes_2/internal/handlers/signup"
	"github.com/golovpeter/clever_notes_2/internal/storage"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func main() {
	url := fmt.Sprintf("postgres://%s:%d@%s:%d/%s", storage.User, storage.Password,
		storage.Host, storage.Port, storage.Dbname)
	db, err := sqlx.Connect("pgx", url)

	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/signup", &signup.DbData{Db: db})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
