package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Reewd/WASAproject/service/api/dto"
	"github.com/Reewd/WASAproject/service/api/helpers"
	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// TODO: MIME type validation for uploaded images
func (rt *_router) uploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Handle image upload logic here
	// This is a placeholder function; actual implementation will depend on how images are uploaded

	err := r.ParseMultipartForm(1024)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to parse multipart form")
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
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to save uploaded file")
		return
	}

	if err := rt.db.InsertImage(uuid, filePath); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to store image path in database")
		os.Remove(filePath)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dto.Photo{
		PhotoId: uuid,
		Path:    filePath,
	}); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}
