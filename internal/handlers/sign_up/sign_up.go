package sign_up

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/golovpeter/clever_notes_2/internal/common/hasher"
	"github.com/golovpeter/clever_notes_2/internal/common/make_error_response"
)

type signUpHandler struct {
	Db Database
}

func NewSignUpHandler(db Database) *signUpHandler {
	return &signUpHandler{Db: db}
}

func (s *signUpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Unsupported method",
		})
		return
	}

	defer r.Body.Close()

	var in SignUpIn

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

	if !validateIn(in) {
		w.WriteHeader(http.StatusBadRequest)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Incorrect data input",
		})
		return
	}

	elementExist := false
	err = s.Db.Get(&elementExist, "select exists(select username from users where username = $1)", in.Username)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if elementExist {
		w.WriteHeader(http.StatusBadRequest)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "User already registered!",
		})
		return
	}

	_, err = s.Db.Exec("insert into users (username, password) values ($1, $2)",
		in.Username, hasher.GeneratePasswordHash(in.Password))

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
		ErrorCode:    "0",
		ErrorMessage: "Registration was successful!",
	})

	return
}

func validateIn(in SignUpIn) bool {
	return in.Username != "" && in.Password != ""
}
