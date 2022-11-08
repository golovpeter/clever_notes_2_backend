package make_error_response

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func MakeErrorResponse(w http.ResponseWriter, els ErrorMessage) {
	out, _ := json.Marshal(els)

	wrote, err := w.Write(out)

	if err != nil || wrote != len(out) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
