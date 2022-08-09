package signup

import (
	"clever_notes_2/internal/storage"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func SingUp(w http.ResponseWriter, r *http.Request) {
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

	url := fmt.Sprintf("postgres://%s:%d@%s:%d/%s", storage.User, storage.Password,
		storage.Host, storage.Port, storage.Dbname)
	db, err := sqlx.Connect("pgx", url)

	if err != nil {
		_, _ = fmt.Fprint(w, "The connection to the database is not established")
		log.Fatalln(err)
		return
	}

	defer db.Close()

	var elementExist bool
	err = db.Get(&elementExist, "select exists(select email from users where email = $1)", in.Email)

	if err != nil {
		log.Fatalln(err)
		return
	}

	if in.Password == in.ConformPass {
		if !elementExist {

			tx := db.MustBegin()

			tx.MustExec("insert into users (email, password, conform_pass) values ($1, $2, $3)",
				in.Email, in.Password, in.ConformPass)

			err = tx.Commit()

			if err != nil {
				log.Fatalln(err)
				return
			}

			_, _ = fmt.Fprintf(w, "User succesful register")

		} else {
			_, _ = fmt.Fprint(w, "Element already registered")
		}
	} else {
		_, _ = fmt.Fprint(w, "Passwords don't match")
	}

}
