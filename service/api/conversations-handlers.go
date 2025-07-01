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

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the conversation details
	var req dto.CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	username, err := rt.db.GetUsername(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get username")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Participants = append(req.Participants, username) // Add the current user to participants

	if req.IsGroup {
		rt.createGroup(w, ctx, req)
	} else {
		rt.createPrivateConversation(w, ctx, req)
	}
}

func (rt *_router) createGroup(w http.ResponseWriter, ctx reqcontext.RequestContext, req dto.CreateConversationRequest) {
	// Validate input
	if len(req.Name) < constraints.MinGroupNameLength || len(req.Name) > constraints.MaxGroupNameLength {
		http.Error(w, fmt.Sprintf("Group name must be between %d and %d characters", constraints.MinGroupNameLength, constraints.MaxGroupNameLength), http.StatusBadRequest)
		return
	}

	if len(req.Participants) < 2 {
		http.Error(w, "A group must have at least 2 participants", http.StatusBadRequest)
		return
	}

	// Extract Photo
	photoId, Photo := helpers.ExtractPhoto(req.Photo)

	// Insert conversation into the database
	conversationId, err := rt.db.InsertConversation(req.Name, req.Participants, req.IsGroup, photoId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to create conversation")
		return
	}

	// Retrieve participants
	database_participants, err := rt.db.GetParticipants(conversationId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve participants")
		return
	}

	participants := helpers.ConvertPublicUsers(database_participants)

	// Respond with the created conversation
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dto.Chat{
		ConversationId: conversationId,
		Name:           req.Name,
		Participants:   participants,
		IsGroup:        req.IsGroup,
		Photo:          Photo,
	})

	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}

func (rt *_router) createPrivateConversation(w http.ResponseWriter, ctx reqcontext.RequestContext, req dto.CreateConversationRequest) {
	if len(req.Participants) != 2 {
		http.Error(w, "A private conversation must have exactly 2 participants", http.StatusBadRequest)
		return
	}

	conversationId, err := rt.db.PrivateConversationExists(req.Participants)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to check for existing private conversation")
		return
	}

	if conversationId == 0 {
		conversationId, err = rt.db.InsertConversation(req.Name, req.Participants, req.IsGroup, nil)
		ctx.Logger.WithField("conversationId", conversationId).Info("Created new private conversation")
		if err != nil {
			helpers.HandleInternalServerError(ctx, w, err, "Failed to create conversation")
			return
		}
	}

	database_participants, err := rt.db.GetParticipants(conversationId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve participants")
		return
	}

	participants := helpers.ConvertPublicUsers(database_participants)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dto.Chat{
		ConversationId: conversationId,
		Name:           req.Name,
		Participants:   participants,
		IsGroup:        req.IsGroup,
	})

	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Retrieve the user's conversations from the database
	databaseConversations, err := rt.db.GetConversationsByUserId(ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve conversations")
		return
	}

	// Pre-allocate conversations slice
	var conversations = make([]dto.ConversationPreview, 0, len(databaseConversations))
	for _, dbConv := range databaseConversations {
		databaseLastMessage, err := rt.db.GetLastMessage(dbConv.ConversationId)
		if err != nil {
			helpers.HandleInternalServerError(ctx, w, err, fmt.Sprintf("Failed to retrieve last message for conversation %d", dbConv.ConversationId))
			return
		}

		var lastMessage *dto.SentMessage
		if databaseLastMessage != nil {
			msg := helpers.ConvertToSentMessage(*databaseLastMessage)
			lastMessage = &msg
		}

		conversations = append(conversations, dto.ConversationPreview{
			ConversationId: dbConv.ConversationId,
			Name:           dbConv.Name,
			Participants:   helpers.ConvertPublicUsers(dbConv.Participants),
			IsGroup:        dbConv.IsGroup,
			Photo:          helpers.ConvertPhoto(dbConv.Photo),
			LastMessage:    lastMessage,
		})
	}

	resp := map[string]interface{}{
		"conversations": conversations, // Wrap the conversations array with a key
	}

	err = rt.db.InsertDelivered(ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to insert delivered status")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}

}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationIdPath := ps.ByName("conversationId")
	if conversationIdPath == "" {
		ctx.Logger.Error("Debug: conversationIdPath is empty")
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	conversationId, err := strconv.ParseInt(conversationIdPath, 10, 64) // Ensure conversationId is a valid integer
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to parse conversationIdPath")
		return
	}

	// Retrieve the conversation from the database
	database_conversation, err := rt.db.GetConversationById(conversationId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve conversation")
		return
	}

	isIn := false
	username, err := rt.db.GetUsername(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get username")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for _, participant := range database_conversation.Participants {
		if participant.Username == username {
			isIn = true
			break
		}
	}

	if !isIn {
		ctx.Logger.Error("User is not a participant of the conversation")
		http.Error(w, "You are not a participant of this conversation", http.StatusForbidden)
		return
	}

	if database_conversation == nil {
		ctx.Logger.Error("Debug: database_conversation is nil")
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	database_chat, err := rt.db.GetChat(conversationId)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to retrieve chat messages")
		return
	}

	err = rt.db.InsertRead(conversationId, ctx.UserID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to insert read status")
		return
	}

	messages := helpers.ConvertToSentMessages(database_chat)
	participants := helpers.ConvertPublicUsers(database_conversation.Participants)
	name := database_conversation.Name
	isGroup := database_conversation.IsGroup
	photo := helpers.ConvertPhoto(database_conversation.Photo)

	conversation := dto.Chat{
		ConversationId: conversationId,
		Name:           name,
		Participants:   participants,
		IsGroup:        isGroup,
		Photo:          photo,
		Messages:       messages,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(conversation)
	if err != nil {
		helpers.HandleInternalServerError(ctx, w, err, "Failed to encode JSON response")
		return
	}
}
