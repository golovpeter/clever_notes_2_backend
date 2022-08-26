package sign_up

import (
	"encoding/json"
	"github.com/golovpeter/clever_notes_2/internal/common/hasher"
	"github.com/golovpeter/clever_notes_2/internal/common/make_response"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type signUpHandler struct {
	Db *sqlx.DB
}

func NewSignUpHandler(db *sqlx.DB) *signUpHandler {
	return &signUpHandler{Db: db}
}

func (s *signUpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "Unsupported method",
		})
		return
	}

	defer r.Body.Close()

	var in SignUpIn

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

	var elementExist bool
	err = s.Db.Get(&elementExist, "select exists(select username from users where username = $1)", in.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	if elementExist {
		w.WriteHeader(http.StatusBadRequest)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "User already registered!",
		})
		return
	}

	_, err = s.Db.Query("insert into users (username, password) values ($1, $2)",
		in.Username, hasher.GeneratePasswordHash(in.Password))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	make_response.MakeResponse(w, map[string]string{
		"errorCode": "0",
		"message":   "Registration was successful!",
	})

	return
}

func validateIn(in SignUpIn) bool {
	return in.Username != "" && in.Password != ""
}
