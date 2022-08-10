package signin

import (
	"encoding/json"
	"fmt"
	"github.com/golovpeter/clever_notes_2/internal/common/hasher"
	"github.com/golovpeter/clever_notes_2/internal/common/token_generator"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type signInHandler struct {
	Db *sqlx.DB
}

func NewSignInHandler(db *sqlx.DB) *signInHandler {
	return &signInHandler{Db: db}
}

func (s *signInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = fmt.Fprint(w, "Unsupported method")
		return
	}

	defer r.Body.Close()

	var in SignIn

	err := json.NewDecoder(r.Body).Decode(&in)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Incorrect data input")
		return
	}

	var elementExist bool
	err = s.Db.Get(&elementExist, "select exists(select username from users where username = $1)", in.Username)

	if err != nil {
		log.Fatalln()
		return
	}

	if elementExist {
		var userData User
		err = s.Db.Get(&userData, "select user_id, username, password from users where  username = $1", in.Username)

		if in.Username == userData.Username && hasher.GeneratePasswordHash(in.Password) == userData.Password {
			token, err := token_generator.GenerateJWT(in.Username)

			if err != nil {
				log.Fatalln(err)
				return
			}

			out, err := json.Marshal(map[string]string{"token": token})

			if err != nil {
				log.Fatalln(err)
				return
			}

			tx := s.Db.MustBegin()

			tx.MustExec("insert into tokens values ((select user_id from users where users.user_id = $1), $2)", userData.User_id, token)

			tx.Commit()

			wrote, err := w.Write(out)
			if err != nil || wrote != len(out) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprint(w, "Invalid username or password")
		}
	} else {
		_, _ = fmt.Fprint(w, "The user is not registered")
	}
}
