package make_error_response

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMessage struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func MakeErrorResponse(w http.ResponseWriter, els ErrorMessage) {

	out, err := json.Marshal(els)

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
