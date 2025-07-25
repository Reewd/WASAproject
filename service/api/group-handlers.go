package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Reewd/WASAproject/service/api/constraints"
	"github.com/Reewd/WASAproject/service/api/dto"
	"github.com/Reewd/WASAproject/service/api/helpers"
	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.AddToGroupRequest

	conversationIdPath := ps.ByName("conversationId")
	if conversationIdPath == "" {
		ctx.Logger.Error("Debug: conversationIdPath is empty")
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	conversationId, err := strconv.ParseInt(conversationIdPath, 10, 64) // Ensure conversationId is a valid integer
	if err != nil {
		ctx.Logger.WithError(err).Error("Debug: Failed to parse conversationIdPath")
		http.Error(w, "The ID should be an integer", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(conversationId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check participant existence")
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check participant existence")
		return
	}

	if !exists {
		ctx.Logger.Error("User is not a participant of the conversation")
		http.Error(w, "You are not a participant of this conversation", http.StatusForbidden)
		return
	}

	if len(req.Participants) == 0 {
		ctx.Logger.Error("No participants specified in the request")
		http.Error(w, "At least one participant must be specified", http.StatusBadRequest)
		return
	}

	participantsIds, err := rt.db.GetUsersIds(req.Participants)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get user IDs")
		helpers.HandleInternalServerError(ctx, w, err, "Failed to get user IDs")
		return
	}

	err = rt.db.InsertParticipants(conversationId, participantsIds)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to add participants to group")
		helpers.HandleInternalServerError(ctx, w, err, "Failed to add participants to group")
		return
	}

	participants, err := rt.db.GetParticipants(conversationId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve participants")
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve participants")
		return
	}

	resp := helpers.ConvertUsers(participants)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode JSON response")
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.LeaveGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(req.ConversationId, ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check participant existence")
		return
	}

	if !exists {
		http.Error(w, "You are not a participant of this conversation", http.StatusForbidden)
		return
	}

	err = rt.db.RemoveParticipant(req.ConversationId, ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to leave group")
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content response for successful leave
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.SetGroupNameRequest

	conversationIdPath := ps.ByName("conversationId")
	if conversationIdPath == "" {
		ctx.Logger.Error("Debug: conversationIdPath is empty")
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	conversationId, err := strconv.ParseInt(conversationIdPath, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Debug: Failed to parse conversationIdPath")
		http.Error(w, "The ID should be an integer", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(conversationId, ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check participant existence")
		return
	}

	if !exists {
		http.Error(w, "You are not a participant of this conversation", http.StatusForbidden)
		return
	}

	if len(req.Name) < constraints.MinGroupNameLength || len(req.Name) > constraints.MaxGroupNameLength {
		http.Error(w, fmt.Sprintf("Group name must be between %d and %d characters", constraints.MinGroupNameLength, constraints.MaxGroupNameLength), http.StatusBadRequest)
		return
	}

	err = rt.db.UpdateGroupName(conversationId, req.Name)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to update group name")
		return
	}
	resp := map[string]string{"name": req.Name}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var req dto.SetGroupPhotoRequest

	conversationIdPath := ps.ByName("conversationId")
	if conversationIdPath == "" {
		ctx.Logger.Error("Debug: conversationIdPath is empty")
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	conversationId, err := strconv.ParseInt(conversationIdPath, 10, 64) // Ensure conversationId is a valid integer
	if err != nil {
		ctx.Logger.WithError(err).Error("Debug: Failed to parse conversationIdPath")
		http.Error(w, "The ID should be an integer", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(conversationId, ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check participant existence")
		return
	}
	if !exists {
		http.Error(w, "You are not a participant of this conversation", http.StatusForbidden)
		return
	}

	err = rt.db.UpdateGroupPhoto(conversationId, req.Photo.PhotoId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to update group photo")
		return
	}
	var resp dto.Photo
	resp.PhotoId = req.Photo.PhotoId
	resp.Path = req.Photo.Path
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}
