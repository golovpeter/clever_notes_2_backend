package sign_in

import (
	"encoding/json"
	"github.com/golovpeter/clever_notes_2/internal/common/hasher"
	"github.com/golovpeter/clever_notes_2/internal/common/make_error_response"
	"github.com/golovpeter/clever_notes_2/internal/common/token_generator"
	"log"
	"net/http"
)

var generateJWT = func(username string) (string, error) {
	jwt, err := token_generator.GenerateJWT(username)
	return jwt, err
}

var generateRefreshJWT = func() (string, error) {
	jwt, err := token_generator.GenerateRefreshJWT()
	return jwt, err
}

type signInHandler struct {
	Db Database
}

func NewSignInHandler(db Database) *signInHandler {
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

	elementExist := false
	err = s.Db.Get(&elementExist, "select exists(select username from users where username = $1)", in.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println()
		return
	}

	if !elementExist {
		w.WriteHeader(http.StatusBadRequest)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "The user is not registered!",
		})
		return
	}

	var userData User
	err = s.Db.Get(&userData, "select user_id, username, password from users where  username = $1", in.Username)

	if in.Username != userData.Username || hasher.GeneratePasswordHash(in.Password) != userData.Password {
		w.WriteHeader(http.StatusInternalServerError)
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "Incorrect password!",
		})
		return
	}

	tokensExist := false
	err = s.Db.Get(&tokensExist, "select exists(select user_id from tokens where user_id = $1)", userData.User_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if tokensExist {
		_, err = s.Db.Exec("delete from tokens where user_id = $1", userData.User_id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	accessToken, _ := generateJWT(in.Username)
	refreshToken, _ := generateRefreshJWT()

	_, err = s.Db.Exec("insert into tokens values ((select user_id from users where users.user_id = $1), $2, $3)",
		userData.User_id,
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

	wrote, err := w.Write(out)

	if err != nil || wrote != len(out) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func validateIn(in SignIn) bool {
	return in.Username != "" && in.Password != ""
}
