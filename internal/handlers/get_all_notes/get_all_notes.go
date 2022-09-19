package get_all_notes

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golovpeter/clever_notes_2/internal/common/make_error_response"
	"github.com/golovpeter/clever_notes_2/internal/common/parse_auth_header"
	"github.com/golovpeter/clever_notes_2/internal/common/token_generator"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type getAllNotesHandel struct {
	db *sqlx.DB
}

type Note struct {
	NoteId  string `json:"note_id" db:"note_id"`
	Caption string `json:"note_caption" db:"note_caption"`
	Text    string `json:"note" db:"note"`
}

func NewGetAllNotesHandler(db *sqlx.DB) *getAllNotesHandel {
	return &getAllNotesHandel{db: db}
}

func (g *getAllNotesHandel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Unsupported method",
		})
		return
	}

	accessToken, err := parse_auth_header.ParseAuthHeader(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	tokenExist := false
	err = g.db.Get(&tokenExist, "select exists(select access_token from tokens where access_token = $1)", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !tokenExist {
		w.WriteHeader(http.StatusUnauthorized)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "The user is not authorized",
		})
		return
	}

	err = token_generator.ValidateToken(accessToken)

	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		w.WriteHeader(http.StatusUnauthorized)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "access token expired",
		})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var userId int
	err = g.db.Get(&userId, "select user_id from tokens where access_token = $1", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if userId == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "there are no such user",
		})
		return
	}

	notes := make([]Note, 0)

	err = g.db.Select(&notes, "select note_id, note_caption, note from notes where user_id = $1", userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	out, err := json.Marshal(GetAllNotesOut{
		Notes: notes,
	})

	wrote, err := w.Write(out)
	if err != nil || wrote != len(out) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
