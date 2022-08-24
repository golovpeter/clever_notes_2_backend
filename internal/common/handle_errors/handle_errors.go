package handle_errors

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, errorMessage string) {
	out, err := json.Marshal(map[string]string{"errorCode": "1", "errorMessage": errorMessage})

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
