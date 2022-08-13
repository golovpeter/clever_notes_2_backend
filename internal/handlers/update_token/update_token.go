package update_token

import (
	"encoding/json"
	"fmt"
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

func (u *updateTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = fmt.Fprint(w, "Unsupported method")
		return
	}

	defer r.Body.Close()

	var in UpdateTokenIn

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

	_, err = token_generator.ParseToken(in.AccessToken)

	if err != nil && err.Error() == "Token is expired" {

		var tokenExist bool
		err := u.db.Get(&tokenExist,
			"select exists(select access_token, refresh_token from tokens where access_token = $1 and  refresh_token = $2)",
			in.AccessToken,
			in.RefreshToken)

		if err != nil {
			log.Fatalln(err)
			return
		}

		if tokenExist {

			var username string
			err := u.db.Get(&username,
				"select username from users inner join tokens n on users.user_id = n.user_id where refresh_token = $1", in.RefreshToken)

			if err != nil {
				log.Fatalln(err)
			}

			newAccessToken, err := token_generator.GenerateJWT(username)

			if err != nil {
				log.Fatalln(err)
				return
			}

			newRefreshToken, err := token_generator.GenerateRefreshJWT()

			_, err = u.db.Query("update tokens set access_token = $1, refresh_token = $2 where refresh_token = $3",
				newAccessToken,
				newRefreshToken,
				in.RefreshToken)

			if err != nil {
				log.Fatalln(err)
				return
			}

			out, err := json.Marshal(map[string]string{"access_token": newAccessToken, "refresh_token": newRefreshToken})

			wrote, err := w.Write(out)
			if err != nil || wrote != len(out) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, "There are no such tokens")
		}
	} else {
		_, _ = fmt.Fprint(w, "The token has not expired yet")
	}
}

func validateIn(in UpdateTokenIn) bool {
	return in.AccessToken != "" && in.RefreshToken != ""
}
