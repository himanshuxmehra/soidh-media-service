package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	ffmpeg "soidh-media-service/internal/ffmpeg"
)

func (h *Handlers) UploadFile(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountId")
	folderID := chi.URLParam(r, "folderId")
	media_id := r.FormValue("mediaID")

	if accountID == "" || folderID == "" || media_id == "" {
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
	filePath := filepath.Join(uploadDir, media_id+filepath.Ext(header.Filename))
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
	media_id_db, err := h.saveMediaInfo(accountID, folderID, media_id, filePath)
	if err != nil {
		h.logger.Error("Failed to save media info", zap.Error(err))
		http.Error(w, "Failed to save file info", http.StatusInternalServerError)
		return
	}
	fmt.Println(header)
	// Return success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File uploaded successfully. Media ID: %s", media_id_db)
}

func (h *Handlers) saveMediaInfo(accountID, folderID, media_id, filePath string) (string, error) {
	query := `INSERT INTO media (account_id, folder_id, media_id, media_type) VALUES ($1, $2, $3, 'image') RETURNING media_id`
	var media_id_db string
	err := h.db.QueryRow(query, accountID, folderID, media_id).Scan(&media_id_db)
	if err != nil {
		return "", err
	}
	return media_id_db, nil
}

func (h *Handlers) UploadVideo(w http.ResponseWriter, r *http.Request){
	media_id := r.FormValue("mediaID")
	account_id := chi.URLParam(r, "accountId")
	folder_id := chi.URLParam(r, "folderId")

	if media_id == "" && account_id == "" && folder_id == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}
	
	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get video file from form", zap.Error(err))
		http.Error(w, "Failed to get video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := filepath.Join("uploads", account_id, folder_id, media_id)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err!=nil{
		h.logger.Error("Failed to create upload directory", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(uploadDir, media_id+filepath.Ext(header.Filename))
	dst, err := os.Create(filePath)
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		h.logger.Error("Failed to copy file", zap.Error(err))
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Save file info to database
	media_id_db, err := h.saveVideoInfo(account_id, folder_id, media_id, filePath)
	if err != nil {
		h.logger.Error("Failed to save media info", zap.Error(err))
		http.Error(w, "Failed to save file info", http.StatusInternalServerError)
		return
	}

	ffmpeg.VideoConversion(filepath.Join("uploads", account_id, folder_id, media_id, filepath.Ext(header.Filename)))

	// Return success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File uploaded successfully. Media ID: %s", media_id_db)
}


func (h *Handlers) saveVideoInfo(accountID, folderID, media_id, filePath string) (string, error) {
	query := `INSERT INTO media (account_id, folder_id, media_id, media_type) VALUES ($1, $2, $3, 'video') RETURNING media_id`
	var media_id_db string
	err := h.db.QueryRow(query, accountID, folderID, media_id).Scan(&media_id_db)
	if err != nil {
		return "", err
	}
	return media_id_db, nil
}