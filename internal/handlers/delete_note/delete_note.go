package delete_note

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golovpeter/clever_notes_2/internal/common/make_error_response"
	"github.com/golovpeter/clever_notes_2/internal/common/parse_auth_header"
	"github.com/golovpeter/clever_notes_2/internal/common/token_generator"
	"log"
	"net/http"
)

var parseAuthHeader = func(w http.ResponseWriter, r *http.Request) (string, error) {
	accessToken, err := parse_auth_header.ParseAuthHeader(w, r)
	return accessToken, err
}

var validateToken = func(accessToken string) error {
	err := token_generator.ValidateToken(accessToken)
	return err
}

type deleteNoteHandler struct {
	db Database
}

func NewDeleteNoteHandler(db Database) *deleteNoteHandler {
	return &deleteNoteHandler{db: db}
}

func (d *deleteNoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Unsupported method",
		})
		return
	}

	var in DeleteNoteIn

	err := json.NewDecoder(r.Body).Decode(&in)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Incorrect data input",
		})
		return
	}

	if !validateIn(in) {
		w.WriteHeader(http.StatusBadRequest)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Incorrect data input",
		})
		return
	}

	accessToken, err := parseAuthHeader(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	tokenExist := false
	err = d.db.Get(&tokenExist, "select exists(select access_token from tokens where access_token = $1)", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !tokenExist {
		w.WriteHeader(http.StatusInternalServerError)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "There are no such tokens",
		})
		return
	}

	err = validateToken(accessToken)

	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		w.WriteHeader(http.StatusUnauthorized)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Access token expired",
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

	if userId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "there is no such user",
		})
		return
	}

	var noteUserId int
	err = d.db.Get(&noteUserId, "select user_id from notes where note_id = $1", in.NoteId)

	if noteUserId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "there is no such note",
		})
		return
	}

	if noteUserId != userId {
		w.WriteHeader(http.StatusBadRequest)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "This note belongs to another user",
		})
		return
	}

	_, err = d.db.Exec("delete from notes where note_id = $1", in.NoteId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
		ErrorCode:    "0",
		ErrorMessage: "note successful deleted",
	})
}

func validateIn(in DeleteNoteIn) bool {
	return in.NoteId != 0
}
