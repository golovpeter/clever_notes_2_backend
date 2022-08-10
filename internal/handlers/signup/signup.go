package signup

import (
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
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

	var in SignUp

	err := json.NewDecoder(r.Body).Decode(&in)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Incorrect data input")
		return
	}

	if err != nil {
		_, _ = fmt.Fprint(w, "The connection to the database is not established")
		log.Fatalln(err)
		return
	}

	var elementExist bool
	err = s.Db.Get(&elementExist, "select exists(select username from users where username = $1)", in.Username)

	if err != nil {
		log.Fatalln(err)
		return
	}

	if !elementExist {

		tx := s.Db.MustBegin()

		tx.MustExec("insert into users (username, password) values ($1, $2)",
			in.Username, in.Password)

		err = tx.Commit()

		if err != nil {
			log.Fatalln(err)
			return
		}

		_, _ = fmt.Fprintf(w, "User succesful register")
		return

	} else {
		_, _ = fmt.Fprint(w, "Element already registered")
		return
	}
}
