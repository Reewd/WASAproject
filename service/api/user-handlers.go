package api

import (
	"encoding/json"
	"net/http"

	"github.com/Reewd/WASAproject/service/api/dto"

	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the username
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate username
	if req.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	// Get or create user ID
	id, err := rt.db.Login(req.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Login failed")
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// Return the user ID to be used as bearer token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the new username
	var req dto.SetUsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate username
	if req.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	// Update the username in the database
	if err := rt.db.UpdateUsername(req.Username, ctx.UserID); err != nil {
		ctx.Logger.WithError(err).Error("Failed to set username")
		http.Error(w, "Failed to set username", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content response for successful update
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the new photo ID

	var req dto.PhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate photo ID
	if req.PhotoId == "" {
		http.Error(w, "Photo ID cannot be empty", http.StatusBadRequest)
		return
	}

	// Update the photo ID in the database
	if err := rt.db.UpdateUserPhoto(req.PhotoId, ctx.UserID); err != nil {
		ctx.Logger.WithError(err).Error("Failed to set photo")
		http.Error(w, "Failed to set photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content response for successful update
}
