package add_note

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

type addNoteHandler struct {
	db *sqlx.DB
}

func NewAddNoteHandler(db *sqlx.DB) *addNoteHandler {
	return &addNoteHandler{db: db}
}

func (a *addNoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "Unsupported method",
		})
		return
	}

	var in AddNoteIn

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
		w.WriteHeader(http.StatusUnauthorized)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "Incorrect header input",
		})
		return
	}

	var tokenExist bool
	err = a.db.Get(&tokenExist, "select exists(select access_token from tokens where access_token = $1)", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !tokenExist {
		w.WriteHeader(http.StatusUnauthorized)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "The user is not authorized!",
		})
		return
	}

	err = token_generator.ValidateToken(accessToken)

	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		w.WriteHeader(http.StatusUnauthorized)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "access token expired",
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

	_, err = a.db.Query("insert into notes(user_id, note_caption, note) values ($1, $2, $3)", userId, in.NoteCaption, in.Note)

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

	out, err := json.Marshal(map[string]int{"note_id": noteId})

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

func validateIn(in AddNoteIn) bool {
	return in.Note != ""
}
