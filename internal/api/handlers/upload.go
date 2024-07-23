package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (h *Handlers) UploadFile(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountId")
	folderID := chi.URLParam(r, "folderId")
	mediaID := r.FormValue("mediaID")

	if accountID == "" || folderID == "" || mediaID == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get file from form", zap.Error(err))
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the upload directory if it doesn't exist
	uploadDir := filepath.Join("uploads", accountID, folderID)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		h.logger.Error("Failed to create upload directory", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create the file
	filePath := filepath.Join(uploadDir, header.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		h.logger.Error("Failed to create file", zap.Error(err))
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the local file system
	if _, err := io.Copy(dst, file); err != nil {
		h.logger.Error("Failed to copy file", zap.Error(err))
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Save file info to database
	media_ID, err := h.saveMediaInfo(accountID, folderID, mediaID, filePath)
	if err != nil {
		h.logger.Error("Failed to save media info", zap.Error(err))
		http.Error(w, "Failed to save file info", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File uploaded successfully. Media ID: %s", media_ID)
}

func (h *Handlers) saveMediaInfo(accountID, folderID, mediaID, filePath string) (string, error) {
	query := `INSERT INTO media (account_id, folder_id, media_id) VALUES ($1, $2, $3) RETURNING id`
	var media_ID string
	err := h.db.QueryRow(query, accountID, folderID, mediaID).Scan(&mediaID)
	if err != nil {
		return "", err
	}
	return media_ID, nil
}
