package update_note

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golovpeter/clever_notes_2/internal/common/token_generator"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type updateNoteHandler struct {
	db *sqlx.DB
}

func NewUpdateNoteHandler(db *sqlx.DB) *updateNoteHandler {
	return &updateNoteHandler{db: db}
}

func (u *updateNoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = fmt.Fprint(w, "Unsupported method")
		return
	}

	var in UpdateNoteIn

	err := json.NewDecoder(r.Body).Decode(&in)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Incorrect data input")
		return
	}

	if !validateIn(in) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Incorrect data input")
		return
	}

	accessToken := r.Header.Get("access_token")

	if accessToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Incorrect header input")
		return
	}

	var tokenExist bool
	err = u.db.Get(&tokenExist, "select exists(select access_token from tokens where access_token = $1)", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	if !tokenExist {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "There are no such tokens")
		return
	}

	err = token_generator.ValidateToken(accessToken)

	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "Access token expired")
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	var userId int
	err = u.db.Get(&userId, "select user_id from tokens where access_token = $1", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	_, err = u.db.Query("update notes set note = $1 where note_id = $2 and user_id = $3",
		in.NewNote,
		in.NoteId,
		userId)

	var noteUserId int
	err = u.db.Get(&noteUserId, "select user_id from notes where note_id = $1", in.NoteId)

	if noteUserId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "There is no such note")
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	if noteUserId != userId {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "The user ID does not match")
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	_, _ = fmt.Fprint(w, "Note was updated")
}

func validateIn(in UpdateNoteIn) bool {
	return in.NewNote != "" && in.NoteId != 0
}
