package utils

import (
	"encoding/json"
	"net/http"
)

// RespondJSON creates the HTTP answer with JSON content
// w the HTTP writer used to send the response
// status the HTTP status of answer
// payload the message to send
func RespondJSON(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	//enableCors(&w)
	if (*r).Method == "OPTIONS" {
		return
	}
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}
