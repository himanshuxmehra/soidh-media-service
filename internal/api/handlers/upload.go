package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

type Media struct {
}

func (h *Handlers) Upload(w http.ResponseWriter, r *http.Request) {

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.logger.Info("there was no file in the request")
		return
	}

	defer file.Close()

	os.MkdirAll("uploads", os.ModePerm)

	filepath := filepath.Join("uploads", handler.Filename)

	dst, err := os.Create(filepath)
	if err != nil {
		h.logger.Info("Failed to create file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("File submitted successfully"))

}
