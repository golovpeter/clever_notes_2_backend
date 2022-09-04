package delete_note

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golovpeter/clever_notes_2/internal/common/make_response"
	"github.com/golovpeter/clever_notes_2/internal/common/token_generator"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type deleteNoteHandler struct {
	db *sqlx.DB
}

func NewDeleteNoteHandler(db *sqlx.DB) *deleteNoteHandler {
	return &deleteNoteHandler{db: db}
}

func (d *deleteNoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "Unsupported method",
		})
		return
	}

	var in DeleteNoteIn

	err := json.NewDecoder(r.Body).Decode(&in)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "Incorrect data input",
		})
		return
	}

	if !validateIn(in) {
		w.WriteHeader(http.StatusBadRequest)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "Incorrect data input",
		})
		return
	}

	accessToken := r.Header.Get("access_token")

	if accessToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "Incorrect header input",
		})
		return
	}

	var tokenExist bool
	err = d.db.Get(&tokenExist, "select exists(select access_token from tokens where access_token = $1)", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !tokenExist {
		w.WriteHeader(http.StatusInternalServerError)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "There are no such tokens",
		})
		return
	}

	err = token_generator.ValidateToken(accessToken)

	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		w.WriteHeader(http.StatusUnauthorized)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "Access token expired",
		})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var userId int
	err = d.db.Get(&userId, "select user_id from tokens where access_token = $1", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var noteUserId int
	err = d.db.Get(&noteUserId, "select user_id from notes where note_id = $1", in.NoteId)

	if noteUserId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "There is no such note",
		})
		return
	}

	if noteUserId != userId {
		w.WriteHeader(http.StatusBadRequest)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "This note belongs to another user",
		})
		return
	}

	_, err = d.db.Query("delete from notes where note_id = $1", in.NoteId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	make_response.MakeResponse(w, map[string]string{
		"errorCode":    "0",
		"errorMessage": "note successful deleted",
	})
}

func validateIn(in DeleteNoteIn) bool {
	return in.NoteId != 0
}
