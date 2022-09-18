package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golovpeter/clever_notes_2/internal/common/enable_cors"
	"github.com/golovpeter/clever_notes_2/internal/handlers/add_note"
	"github.com/golovpeter/clever_notes_2/internal/handlers/delete_note"
	"github.com/golovpeter/clever_notes_2/internal/handlers/get_all_notes"
	"github.com/golovpeter/clever_notes_2/internal/handlers/log_out"
	servestatic "github.com/golovpeter/clever_notes_2/internal/handlers/serve_static"
	"github.com/golovpeter/clever_notes_2/internal/handlers/sign_in"
	"github.com/golovpeter/clever_notes_2/internal/handlers/sign_up"
	"github.com/golovpeter/clever_notes_2/internal/handlers/update_note"
	"github.com/golovpeter/clever_notes_2/internal/handlers/update_token"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(1 * time.Second)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println(db.Stats())
			case <-ctx.Done():
				return
			}
		}
	}()

	db.Stats()

	mux := http.NewServeMux()

	// Authentication
	mux.Handle("/sign-up", enable_cors.CORS(sign_up.NewSignUpHandler(db)))
	mux.Handle("/sign-in", enable_cors.CORS(sign_in.NewSignInHandler(db)))
	mux.Handle("/log-out", enable_cors.CORS(log_out.NewLogOutHandler(db)))

	// Working with notes
	mux.Handle("/add-note", enable_cors.CORS(add_note.NewAddNoteHandler(db)))
	mux.Handle("/update-note", enable_cors.CORS(update_note.NewUpdateNoteHandler(db)))
	mux.Handle("/delete-note", enable_cors.CORS(delete_note.NewDeleteNoteHandler(db)))

	mux.Handle("/get-all-notes", enable_cors.CORS(get_all_notes.NewGetAllNotesHandler(db)))
	mux.Handle("/update-token", enable_cors.CORS(update_token.NewUpdateTokenHandler(db)))

	// Serve static content
	mux.Handle("/", enable_cors.CORS(servestatic.NewServeStaticHandler("./static")))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), mux))
}
