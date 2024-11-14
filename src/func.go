package server

import (
	"encoding/json"
	"net/http"
)

func Fetch(url string, w http.ResponseWriter) *http.Response {
	// Make a GET request to the provided URL and handle errors
	data1, err := http.Get(url)
	if err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "failed fetshing data"}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
	}
	//return http.Response
	return data1
}

func DecodeByUs(db *http.Response, pointer any, w http.ResponseWriter) {
	// Decode the JSON response into the provided pointer and handle errors
	if err := json.NewDecoder(db.Body).Decode(&pointer); err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "error decoding"}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
	}
}
