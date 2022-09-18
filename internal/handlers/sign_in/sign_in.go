package sign_in

import (
	"encoding/json"
	"github.com/golovpeter/clever_notes_2/internal/common/hasher"
	"github.com/golovpeter/clever_notes_2/internal/common/make_error_response"
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
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Unsupported method",
		})
		return
	}

	defer r.Body.Close()

	var in SignIn

	err := json.NewDecoder(r.Body).Decode(&in)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	elementExist := []bool{false}
	err = s.Db.Select(&elementExist, "select exists(select username from users where username = $1)", in.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println()
		return
	}

	if !elementExist[0] {
		w.WriteHeader(http.StatusBadRequest)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "The user is not registered!",
		})
		return
	}

	var userData []User
	err = s.Db.Select(&userData, "select user_id, username, password from users where  username = $1", in.Username)

	if in.Username != userData[0].Username || hasher.GeneratePasswordHash(in.Password) != userData[0].Password {
		w.WriteHeader(http.StatusInternalServerError)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Incorrect password!",
		})
		return
	}

	tokensExist := []bool{false}
	err = s.Db.Select(&tokensExist, "select exists(select user_id from tokens where user_id = $1)", userData[0].User_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if tokensExist[0] {
		_, err = s.Db.Query("delete from tokens where user_id = $1", userData[0].User_id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	_, err = s.Db.Exec("delete from tokens where user_id = $1", userData[0].User_id)

	accessToken, err := token_generator.GenerateJWT(in.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	refreshToken, err := token_generator.GenerateRefreshJWT()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	_, err = s.Db.Exec("insert into tokens values ((select user_id from users where users.user_id = $1), $2, $3)",
		userData[0].User_id,
		accessToken,
		refreshToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := SignInOut{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	out, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalln(err)
		return
	}

	wrote, err := w.Write(out)

	if err != nil || wrote != len(out) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func validateIn(in SignIn) bool {
	return in.Username != "" && in.Password != ""
}
