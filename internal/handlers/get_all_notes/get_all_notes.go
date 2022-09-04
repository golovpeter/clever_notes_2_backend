package get_all_notes

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

type getAllNotesHandel struct {
	db *sqlx.DB
}

type Note struct {
	Caption string `json:"note_caption"`
	Text    string `json:"note"`
}

func NewGetAllNotesHandler(db *sqlx.DB) *getAllNotesHandel {
	return &getAllNotesHandel{db: db}
}

func (g *getAllNotesHandel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "Unsupported method",
		})
		return
	}

	//TODO: передавать токен в хедере Authorization
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
	err := g.db.Get(&tokenExist, "select exists(select access_token from tokens where access_token = $1)", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !tokenExist {
		w.WriteHeader(http.StatusUnauthorized)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "The user is not authorized",
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
	err = g.db.Get(&userId, "select user_id from tokens where access_token = $1", accessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	notes := make([]Note, 0)

	rows, err := g.db.Query("select note_caption, note from notes where user_id = $1", userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for rows.Next() {
		var noteCaption, note string
		_ = rows.Scan(&noteCaption, &note)

		el := Note{
			Caption: noteCaption,
			Text:    note,
		}

		notes = append(notes, el)
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
