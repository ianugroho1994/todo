package shared

import (
	"encoding/json"
	"net/http"
)

func WriteMessage(w http.ResponseWriter, status int, msg string) {
	var j struct {
		Msg string `json:"message"`
	}

	j.Msg = msg

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(j)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteMessage(w, status, err.Error())
}
