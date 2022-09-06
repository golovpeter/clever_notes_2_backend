package parse_auth_header

import (
	"errors"
	"github.com/golovpeter/clever_notes_2/internal/common/make_error_response"
	"net/http"
	"strings"
)

func ParseAuthHeader(w http.ResponseWriter, r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")

	headerSplit := strings.Split(header, " ")
	if len(headerSplit) != 2 || headerSplit[0] != "Bearer" {
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "invalid auth header",
		})
		return "", errors.New("invalid auth header")
	}

	if len(headerSplit[1]) == 0 {
		make_error_response.MakeErrorResponse(w, make_error_response.ErrorMessage{
			ErrorCode:    "1",
			ErrorMessage: "token is empty",
		})
		return "", errors.New("token is empty")
	}

	accessToken := headerSplit[1]
	return accessToken, nil
}
