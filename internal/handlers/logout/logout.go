package logout

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type logOutHandler struct {
	Db *sqlx.DB
}

func NewLogOutHandler(db *sqlx.DB) *logOutHandler {
	return &logOutHandler{Db: db}
}

func (l *logOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = fmt.Fprint(w, "Unsupported method")
		return
	}

	defer r.Body.Close()

	var in LogOutIn

	err := json.NewDecoder(r.Body).Decode(&in)

	if err != nil || !validateIn(in) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Incorrect data input")
		return
	}

	if !validateIn(in) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Incorrect data input")
		return
	}

	var tokenExist bool
	err = l.Db.Get(&tokenExist, "select exists(select access_token from tokens where access_token = $1)",
		in.AccessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	if !tokenExist {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "There are no such tokens")
		return
	}

	_, err = l.Db.Query("delete from tokens where access_token = $1", in.AccessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

}

func validateIn(in LogOutIn) bool {
	return in.AccessToken != ""
}
