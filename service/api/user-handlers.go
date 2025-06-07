package api

import (
	"encoding/json"
	"net/http"

	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the username
	var req UsernameRequest
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

func (rt *_router) SetMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the new username
	var req UsernameRequest
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
	if err := rt.db.SetMyUsername(req.Username, ctx.UserID); err != nil {
		ctx.Logger.WithError(err).Error("Failed to set username")
		http.Error(w, "Failed to set username", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content response for successful update
}

func (rt *_router) SetMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the new photo ID

	var req PhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate photo ID
	if req.PhotoID == "" {
		http.Error(w, "Photo ID cannot be empty", http.StatusBadRequest)
		return
	}

	// Update the photo ID in the database
	if err := rt.db.SetMyPhoto(req.PhotoID, ctx.UserID); err != nil {
		ctx.Logger.WithError(err).Error("Failed to set photo")
		http.Error(w, "Failed to set photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content response for successful update
}
