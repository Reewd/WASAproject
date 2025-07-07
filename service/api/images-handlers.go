package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Reewd/WASAproject/service/api/constraints"
	"github.com/Reewd/WASAproject/service/api/dto"
	"github.com/Reewd/WASAproject/service/api/helpers"
	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	err := r.ParseMultipartForm(1024)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to parse multipart form")
		return
	}

	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get image file from form")
		http.Error(w, "Failed to get image file", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			ctx.Logger.WithError(err).Error("Failed to close file")
		}
	}()

	if handler.Size > constraints.MaxFileSize {
		ctx.Logger.Error("File too large")
		http.Error(w, "File is too large. Maximum allowed size is 10MB.", http.StatusBadRequest)
		return
	}

	// Validate MIME type
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to read file for MIME type validation")
		return
	}

	_, err = file.Seek(0, 0) // Reset file pointer after reading

	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to reset file pointer")
		return
	}

	mimeType := http.DetectContentType(buffer)
	isValidMimeType := false
	for _, allowed := range constraints.AllowedMimeTypes {
		if mimeType == allowed {
			isValidMimeType = true
			break
		}
	}

	if !isValidMimeType {
		ctx.Logger.Error("Invalid MIME type: " + mimeType)
		http.Error(w, "Invalid file type.", http.StatusBadRequest)
		return
	}

	uploadDir := "./uploads"
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to create upload directory")
		return
	}

	fileExt := filepath.Ext(handler.Filename)
	uuid := uuid.New().String()
	newFilename := uuid + fileExt
	filePath := filepath.Join(uploadDir, newFilename)
	dst, err := os.Create(filePath)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to create destination file")
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			ctx.Logger.WithError(err).Error("Failed to close destination file")
		}
	}()

	if _, err = io.Copy(dst, file); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to save uploaded file")
		return
	}

	if err := rt.db.InsertImage(uuid, filePath); err != nil {
		err = os.Remove(filePath)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to remove uploaded file after database error")
		}
		helpers.HandleInternalServerError(ctx, w, err, "Failed to store image path in database")
		return
	}

	var resp dto.Photo
	resp.PhotoId = uuid
	resp.Path = filePath

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}
