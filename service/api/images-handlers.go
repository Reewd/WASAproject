package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) UploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Handle image upload logic here
	// This is a placeholder function; actual implementation will depend on how images are uploaded

	err := r.ParseMultipartForm(1024)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to parse multipart form")
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get image file from form")
		http.Error(w, "Failed to get image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "./uploads"
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create upload directory")
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	fileExt := filepath.Ext(handler.Filename)
	uuid := uuid.New().String()
	newFilename := uuid + fileExt
	filePath := filepath.Join(uploadDir, newFilename)
	dst, err := os.Create(filePath)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create destination file")
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		ctx.Logger.WithError(err).Error("Failed to save uploaded file")
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	if err := rt.db.InsertImage(uuid, filePath); err != nil {
		ctx.Logger.WithError(err).Error("Failed to store image path in database")
		// Try to clean up the file if database insertion fails
		os.Remove(filePath)
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"path":   filePath,
	})
}
