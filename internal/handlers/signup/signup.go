package signup

import (
	"encoding/json"
	"fmt"
	"github.com/golovpeter/clever_notes_2/internal/common/hasher"
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
		_, _ = fmt.Fprint(w, "Unsupported method")
		return
	}

	defer r.Body.Close()

	var in SignUpIn

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

	var elementExist bool
	err = s.Db.Get(&elementExist, "select exists(select username from users where username = $1)", in.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	if elementExist {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Element already registered")
		return
	}

	tx := s.Db.MustBegin()

	tx.MustExec("insert into users (username, password) values ($1, $2)",
		in.Username, hasher.GeneratePasswordHash(in.Password))

	err = tx.Commit()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	_, _ = fmt.Fprintf(w, "User succesful register")
	return
}

func validateIn(in SignUpIn) bool {
	return in.Username != "" && in.Password != ""
}
