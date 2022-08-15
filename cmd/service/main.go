package main

import (
	"fmt"
	"github.com/golovpeter/clever_notes_2/internal/handlers/add_note"
	"github.com/golovpeter/clever_notes_2/internal/handlers/log_out"
	"github.com/golovpeter/clever_notes_2/internal/handlers/sign_in"
	"github.com/golovpeter/clever_notes_2/internal/handlers/sign_up"
	"github.com/golovpeter/clever_notes_2/internal/handlers/update_note"
	"github.com/golovpeter/clever_notes_2/internal/handlers/update_token"
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

	// Authentication
	mux.Handle("/sign-up", sign_up.NewSignUpHandler(db))
	mux.Handle("/sign-in", sign_in.NewSignInHandler(db))
	mux.Handle("/logo-ut", log_out.NewLogOutHandler(db))

	// Working with notes
	mux.Handle("/add-note", add_note.NewAddNoteHandler(db))
	mux.Handle("/update-note", update_note.NewUpdateNoteHandler(db))

	mux.Handle("/update-token", update_token.NewUpdateTokenHandler(db))

	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), mux))

	defer db.Close()
}
