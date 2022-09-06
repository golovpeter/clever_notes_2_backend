package log_out

import (
	"github.com/golovpeter/clever_notes_2/internal/common/make_error_response"
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Unsupported method",
		})
		return
	}

	defer r.Body.Close()

	var in LogOutIn

	in.AccessToken = r.Header.Get("access_token")

	if !validateIn(in) {
		w.WriteHeader(http.StatusBadRequest)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Incorrect header input",
		})
		return
	}

	var tokenExist bool
	err := l.Db.Get(&tokenExist, "select exists(select access_token from tokens where access_token = $1)",
		in.AccessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !tokenExist {
		w.WriteHeader(http.StatusInternalServerError)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "there are no such tokens",
		})
		return
	}

	_, err = l.Db.Query("delete from tokens where access_token = $1", in.AccessToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
		ErrorCode:    "0",
		ErrorMessage: "token successful deleted",
	})

}

func validateIn(in LogOutIn) bool {
	return in.AccessToken != ""
}
