package update_token

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golovpeter/clever_notes_2/internal/common/make_response"
	"github.com/golovpeter/clever_notes_2/internal/common/token_generator"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type updateTokenHandler struct {
	db *sqlx.DB
}

func NewUpdateTokenHandler(db *sqlx.DB) *updateTokenHandler {
	return &updateTokenHandler{db: db}
}

//TODO: исправить log.Println на log.error
func (u *updateTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "unsupported method",
		})
		return
	}

	defer r.Body.Close()

	var in UpdateTokenIn

	in.AccessToken = r.Header.Get("access_token")
	in.RefreshToken = r.Header.Get("refresh_token")

	if !validateIn(in) {
		w.WriteHeader(http.StatusBadRequest)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "incorrect header input",
		})
		return
	}

	err := token_generator.ValidateToken(in.RefreshToken)

	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		w.WriteHeader(http.StatusBadRequest)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "refresh token expired",
		})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var tokenExist bool
	err = u.db.Get(&tokenExist,
		"select exists(select access_token, refresh_token from tokens where access_token = $1 and  refresh_token = $2)",
		in.AccessToken,
		in.RefreshToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !tokenExist {
		w.WriteHeader(http.StatusBadRequest)
		make_response.MakeResponse(w, map[string]string{
			"errorCode":    "1",
			"errorMessage": "there are no suck tokens",
		})
		return
	}

	var username string
	err = u.db.Get(&username,
		"select username from users inner join tokens n on users.user_id = n.user_id where refresh_token = $1", in.RefreshToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	newAccessToken, err := token_generator.GenerateJWT(username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	newRefreshToken, err := token_generator.GenerateRefreshJWT()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	_, err = u.db.Query("update tokens set access_token = $1, refresh_token = $2 where refresh_token = $3",
		newAccessToken,
		newRefreshToken,
		in.RefreshToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	out, err := json.Marshal(map[string]string{"access_token": newAccessToken, "refresh_token": newRefreshToken})

	wrote, err := w.Write(out)
	if err != nil || wrote != len(out) {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func validateIn(in UpdateTokenIn) bool {
	return in.AccessToken != "" && in.RefreshToken != ""
}
