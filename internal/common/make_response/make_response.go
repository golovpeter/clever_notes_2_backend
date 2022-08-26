package make_response

import (
	"encoding/json"
	"log"
	"net/http"
)

func MakeResponse(w http.ResponseWriter, els map[string]string) {

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
