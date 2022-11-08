package add_note

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

var validateToken = func(accessToken string) error {
	err := token_generator.ValidateToken(accessToken)
	return err
}

var parseAuthHeader = func(w http.ResponseWriter, r *http.Request) (string, error) {
	accessToken, err := parse_auth_header.ParseAuthHeader(w, r)
	return accessToken, err
}

type addNoteHandler struct {
	db Database
}

func NewAddNoteHandler(db Database) *addNoteHandler {
	return &addNoteHandler{db: db}
}

func (a *addNoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Unsupported method",
		})
		return
	}

	var in AddNoteIn

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

	accessToken, err := parseAuthHeader(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	tokenExist := false
	err = a.db.Get(&tokenExist, "select exists(select access_token from tokens where access_token = $1)", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !tokenExist {
		w.WriteHeader(http.StatusUnauthorized)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "The user is not authorized!",
		})
		return
	}

	err = validateToken(accessToken)

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
	err = a.db.Get(&userId, "select user_id from tokens where access_token = $1", accessToken)

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

	_, err = a.db.Exec("insert into notes(user_id, note_caption, note) values ($1, $2, $3)", userId, in.NoteCaption, in.Note)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var noteId int
	err = a.db.Get(&noteId, "select max(note_id) from notes")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := AddNoteOut{NoteId: noteId}

	out, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	wrote, err := w.Write(out)

	if err != nil || wrote != len(out) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
