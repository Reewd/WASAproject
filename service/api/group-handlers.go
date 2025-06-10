package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Reewd/WASAproject/service/api/constraints"
	"github.com/Reewd/WASAproject/service/api/dto"
	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.AddToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(req.ConversationId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check participant existence")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "You are not a participant of this conversation", http.StatusForbidden)
		return
	}

	if len(req.Participants) == 0 {
		http.Error(w, "At least one participant must be specified", http.StatusBadRequest)
		return
	}

	participantsIds, err := rt.db.GetUsersIds(req.Participants)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get user IDs")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = rt.db.InsertParticipants(req.ConversationId, participantsIds)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to add participants to group")
		http.Error(w, "Failed to add participants", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // No content response for successful addition
	return
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.LeaveGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(req.ConversationId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check participant existence")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "You are not a participant of this conversation", http.StatusForbidden)
		return
	}

	err = rt.db.RemoveParticipant(req.ConversationId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to leave group")
		http.Error(w, "Failed to leave group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content response for successful leave
	return
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.SetGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(req.ConversationId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check participant existence")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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

	err = rt.db.UpdateGroupName(req.ConversationId, req.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to update group name")
		http.Error(w, "Failed to update group name", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content response for successful update
	return
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.SetGroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(req.ConversationId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check participant existence")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "You are not a participant of this conversation", http.StatusForbidden)
		return
	}
	if len(req.PhotoId) != 20 {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}
	err = rt.db.UpdateGroupPhoto(req.ConversationId, req.PhotoId)

	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to update group photo")
		http.Error(w, "Failed to update group photo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // No content response for successful update
	return
}
