package main

import (
	"fmt"
	"github.com/golovpeter/clever_notes_2/internal/handlers/signin"
	"github.com/golovpeter/clever_notes_2/internal/handlers/signup"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

func main() {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB_NAME"))

	db, err := sqlx.Connect("pgx", url)
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/signup", signup.NewSignUpHandler(db))
	mux.Handle("/signin", signin.NewSignInHandler(db))

	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), mux))

	defer db.Close()
}
