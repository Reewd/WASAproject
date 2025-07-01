package api

import (
	"encoding/json"
	"net/http"

	"github.com/Reewd/WASAproject/service/api/dto"
	"github.com/Reewd/WASAproject/service/api/helpers"

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
	dbUser, err := rt.db.Login(req.Username)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Login failed")
		return
	}

	var resp dto.User
	resp.Username = req.Username
	resp.UserId = dbUser.UserId
	resp.Photo = helpers.ConvertPhoto(dbUser.Photo)

	// Return the user ID to be used as bearer token
	w.Header().Set("Content-Type", "application/json")
	// Add error handling for JSON encoding
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
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
	err := rt.db.UpdateUsername(req.Username, ctx.UserID)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
		helpers.HandleInternalServerError(ctx, w, err, "Failed to set username")
		return
	}
	// Get the updated user data
	dbUser, err := rt.db.GetUser(ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve updated user")
		return
	}

	// Return the updated user
	resp := helpers.ConvertUser(*dbUser)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}

}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the new photo ID

	var req dto.PhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate photo ID
	if req.Photo.PhotoId == "" {
		http.Error(w, "Photo ID cannot be empty", http.StatusBadRequest)
		return
	}

	// Update the photo ID in the database
	if err := rt.db.UpdateUserPhoto(req.Photo.PhotoId, ctx.UserID); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to set photo")
		return
	}

	var resp dto.User
	resp.Photo = &dto.Photo{
		PhotoId: req.Photo.PhotoId,
		Path:    req.Photo.Path,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get all users from the database
	users, err := rt.db.GetAllUsers()
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve users")
		return
	}

	// Convert users to DTOs
	dtoUsers := helpers.ConvertUsers(users)

	resp := map[string][]dto.User{
		"users": dtoUsers,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}
