package api

import (
	"encoding/json"
	"net/http"

	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the conversation details
	var req CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.Participants = append(req.Participants, ctx.Username) // Add the current user to participants

	if req.IsGroup {
		rt.createGroup(w, ctx, req)
	} else {
		rt.createPrivateConversation(w, ctx, req)
	}
}

func (rt *_router) createGroup(w http.ResponseWriter, ctx reqcontext.RequestContext, req CreateConversationRequest) {
	// Create the conversation in the database
	if req.Title == "" {
		http.Error(w, "Title of a group cannot be empty", http.StatusBadRequest)
		return
	}

	if len(req.Participants) < 2 {
		http.Error(w, "A group must have at least 2 participants", http.StatusBadRequest)
		return
	}

	conversationId, err := rt.db.InsertConversation(req.Title, req.Participants, req.IsGroup, req.PhotoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create conversation")
		http.Error(w, "Failed to create conversation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ConversationResponse{conversationId, req.Title, req.Participants, req.IsGroup, req.PhotoId})
}

func (rt *_router) createPrivateConversation(w http.ResponseWriter, ctx reqcontext.RequestContext, req CreateConversationRequest) {
	conversationId, err := rt.db.InsertConversation(req.Title, req.Participants, req.IsGroup, req.PhotoId)

	if len(req.Participants) != 2 {
		http.Error(w, "A private conversation must have exactly 2 participants", http.StatusBadRequest)
		return
	}

	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create conversation")
		http.Error(w, "Failed to create conversation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ConversationResponse{conversationId, req.Title, req.Participants, req.IsGroup, req.PhotoId})
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext)
