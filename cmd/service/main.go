package main

import (
	"fmt"
	"github.com/golovpeter/clever_notes_2/internal/common/enable_cors"
	"github.com/golovpeter/clever_notes_2/internal/handlers/add_note"
	"github.com/golovpeter/clever_notes_2/internal/handlers/delete_note"
	"github.com/golovpeter/clever_notes_2/internal/handlers/get_add_notes"
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
	mux.Handle("/sign-up", enable_cors.CORS(sign_up.NewSignUpHandler(db)))
	mux.Handle("/sign-in", enable_cors.CORS(sign_in.NewSignInHandler(db)))
	mux.Handle("/log-out", enable_cors.CORS(log_out.NewLogOutHandler(db)))

	// Working with notes
	mux.Handle("/add-note", enable_cors.CORS(add_note.NewAddNoteHandler(db)))
	mux.Handle("/update-note", enable_cors.CORS(update_note.NewUpdateNoteHandler(db)))
	mux.Handle("/delete-note", enable_cors.CORS(delete_note.NewDeleteNoteHandler(db)))

	mux.Handle("/get-all-notes", enable_cors.CORS(get_add_notes.NewGetAllNotesHandler(db)))
	mux.Handle("/update-token", enable_cors.CORS(update_token.NewUpdateTokenHandler(db)))

	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), mux))

	defer db.Close()
}
