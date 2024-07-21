package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Welcome to the home page!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
