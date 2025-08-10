package presentation

import (
	"bytes"
	"io"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := http.Post("http://localhost:8000/users/", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		http.Error(w, "Failed to make request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write(responseBody)
}
