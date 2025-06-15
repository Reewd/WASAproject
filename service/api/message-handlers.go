package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Reewd/WASAproject/service/api/constraints"
	"github.com/Reewd/WASAproject/service/api/dto"
	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Content) == 0 || len(req.Content) > constraints.MaxMessageLength {
		http.Error(w, fmt.Sprintf("Message content is required and must not exceed %d characters", constraints.MaxMessageLength), http.StatusBadRequest)
		return
	}

	conversationId, err := strconv.ParseInt(ps.ByName("conversationId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(conversationId, ctx.UserID)
	if err != nil {
		http.Error(w, "Failed to check participant existence", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "You are not a participant in this conversation", http.StatusForbidden)
		return
	}

	if err := rt.db.InsertMessage(conversationId, ctx.UserID, req.Content, req.PhotoId, req.ReplyToMessageId); err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	//TODO: Return the message
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId, err := strconv.ParseInt(ps.ByName("messageId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	conversationId, err := strconv.ParseInt(ps.ByName("conversationId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	dbConversationId, err := rt.db.GetConversationIdFromMessageId(messageId)
	if err != nil {
		http.Error(w, "Failed to retrieve conversation ID from message", http.StatusInternalServerError)
		return
	}

	if dbConversationId != conversationId {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	exists, err := rt.db.ParticipantExists(conversationId, ctx.UserID)
	if err != nil {
		http.Error(w, "Failed to check participant existence", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "You are not a participant in this conversation", http.StatusForbidden)
		return
	}

	senderId, err := rt.db.GetSenderId(messageId)
	if err != nil {
		http.Error(w, "Failed to retrieve message sender", http.StatusInternalServerError)
		return
	}

	if senderId != ctx.UserID {
		http.Error(w, "You are not the sender of this message", http.StatusForbidden)
		return
	}

	if err := rt.db.RemoveMessage(messageId); err != nil {
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.ForwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	conversationId, err := strconv.ParseInt(ps.ByName("conversationId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(conversationId, ctx.UserID)
	if err != nil {
		http.Error(w, "Failed to check participant existence", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "You are not a participant in this conversation", http.StatusForbidden)
		return
	}

	fromConversationId, err := rt.db.GetConversationIdFromMessageId(req.MessageId)
	if err != nil {
		http.Error(w, "Failed to retrieve conversation ID from message", http.StatusInternalServerError)
		return
	}

	exists, err = rt.db.ParticipantExists(fromConversationId, ctx.UserID)
	if err != nil {
		http.Error(w, "Failed to check participant existence in source conversation", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "You are not a participant in the source conversation", http.StatusForbidden)
		return
	}

	if err := rt.db.ForwardMessage(req.MessageId, conversationId, ctx.UserID); err != nil {
		http.Error(w, "Failed to forward message", http.StatusInternalServerError)
		return
	}

	//TODO: Return the forwarded message
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.ReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	messageId, err := strconv.ParseInt(ps.ByName("messageId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	conversationId, err := strconv.ParseInt(ps.ByName("conversationId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}
	exists, err := rt.db.ParticipantExists(conversationId, ctx.UserID)
	if err != nil {
		http.Error(w, "Failed to check participant existence", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "You are not a participant in this conversation", http.StatusForbidden)
		return
	}
	//TODO: Validate req.Content, e.g., check length or allowed characters
	err = rt.db.InsertReaction(messageId, ctx.UserID, req.Content)
	if err != nil {
		http.Error(w, "Failed to comment on message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId, err := strconv.ParseInt(ps.ByName("messageId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	conversationId, err := strconv.ParseInt(ps.ByName("conversationId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}
	exists, err := rt.db.ParticipantExists(conversationId, ctx.UserID)
	if err != nil {
		http.Error(w, "Failed to check participant existence", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "You are not a participant in this conversation", http.StatusForbidden)
		return
	}

	err = rt.db.RemoveReaction(messageId, ctx.UserID)
	if err != nil {
		http.Error(w, "Failed to remove comment from message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
