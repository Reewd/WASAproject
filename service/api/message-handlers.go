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

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Text != nil && req.Photo != nil {
		http.Error(w, "You must send either a message or a photo, or both", http.StatusBadRequest)
	}

	if req.Text != nil {
		if len(*req.Text) == 0 || len(*req.Text) > constraints.MaxMessageLength {
			http.Error(w, fmt.Sprintf("Message content cannot be empty and must not exceed %d characters", constraints.MaxMessageLength), http.StatusBadRequest)
			return
		}
	}

	conversationId, err := strconv.ParseInt(ps.ByName("conversationId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.ParticipantExists(conversationId, ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check participant existence")
		return
	}
	if !exists {
		http.Error(w, "You are not a participant in this conversation", http.StatusForbidden)
		return
	}

	messageId, timestamp, err := rt.db.InsertMessage(conversationId, ctx.UserID, req.Text, &req.Photo.PhotoId, req.ReplyToMessageId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to insert message")
		return
	}

	participantIds, err := rt.db.GetParticipantIds(conversationId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve participant IDs")
		return
	}

	if err := rt.db.InsertSent(messageId, conversationId, participantIds); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to insert sent status")
		return
	}

	dbUser, err := rt.db.GetPublicUser(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get sender's public user information")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var resp dto.SentMessage
	resp.MessageId = messageId
	resp.ConversationId = conversationId
	resp.Timestamp = timestamp
	resp.Photo = &dto.Photo{
		PhotoId: req.Photo.PhotoId,
		Path:    req.Photo.Path,
	}
	resp.SentBy = helpers.ConvertPublicUser(*dbUser)
	resp.Text = req.Text
	resp.ReplyToMessageId = req.ReplyToMessageId
	resp.Status = "sent" // Initial status is "sent"

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
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
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve conversation ID from message")
		return
	}

	if dbConversationId != conversationId {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	exists, err := rt.db.ParticipantExists(conversationId, ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check participant existence")
		return
	}

	if !exists {
		http.Error(w, "You are not a participant in this conversation", http.StatusForbidden)
		return
	}

	senderId, err := rt.db.GetSenderId(messageId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve sender ID from message")
		return
	}

	if senderId != ctx.UserID {
		http.Error(w, "You are not the sender of this message", http.StatusForbidden)
		return
	}

	if err := rt.db.RemoveMessage(messageId); err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to delete message")
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
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check participant existence")
		return
	}
	if !exists {
		http.Error(w, "You are not a participant in this conversation", http.StatusForbidden)
		return
	}

	fromConversationId, err := rt.db.GetConversationIdFromMessageId(req.MessageId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve conversation ID from message")
		return
	}

	exists, err = rt.db.ParticipantExists(fromConversationId, ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check participant existence in source conversation")
		return
	}
	if !exists {
		http.Error(w, "You are not a participant in the source conversation", http.StatusForbidden)
		return
	}

	messageId, timestamp, text, photoId, err := rt.db.ForwardMessage(req.MessageId, conversationId, ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to forward message")
		return
	}

	dbUser, err := rt.db.GetPublicUser(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get sender's public user information")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var respPhoto *dto.Photo
	if photoId == nil {
		path, err := rt.db.GetImagePath(*photoId)
		if err != nil {
			helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve image path")
			return
		}
		respPhoto.PhotoId = *photoId
		respPhoto.Path = path
	}

	var resp dto.SentMessage
	resp.MessageId = messageId
	resp.Timestamp = timestamp
	resp.Photo = respPhoto
	resp.SentBy = helpers.ConvertPublicUser(*dbUser)
	resp.Text = text
	resp.ReplyToMessageId = nil
	resp.Status = "sent" // Initial status is "sent"

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req dto.ReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := helpers.IsSingleEmoji(req.Content)
	if err != nil {
		http.Error(w, "Invalid reaction content: "+err.Error(), http.StatusBadRequest)
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
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check participant existence")
		return
	}
	if !exists {
		http.Error(w, "You are not a participant in this conversation", http.StatusForbidden)
		return
	}

	err = rt.db.InsertReaction(messageId, ctx.UserID, req.Content)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to comment on message")
		return
	}

	reactions, err := rt.db.GetReactions(messageId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve reactions for message")
		return
	}

	resp := helpers.ConvertReactions(reactions)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
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
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check participant existence")
		return
	}
	if !exists {
		http.Error(w, "You are not a participant in this conversation", http.StatusForbidden)
		return
	}

	err = rt.db.RemoveReaction(messageId, ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to remove comment from message")
		return
	}

	reactions, err := rt.db.GetReactions(messageId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve reactions for message")
		return
	}

	resp := helpers.ConvertReactions(reactions)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}
