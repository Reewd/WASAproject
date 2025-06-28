package dto

type LoginRequest struct {
	Username string `json:"username"`
}

type SetUsernameRequest struct {
	Username string `json:"username"`
}

type PhotoRequest struct {
	PhotoId string `json:"photoId"`
}

type CreateConversationRequest struct {
	Name         string   `json:"name,omitempty"`
	Participants []string `json:"participants"`
	IsGroup      bool     `json:"isGroup"`
	PhotoId      *string  `json:"photoId,omitempty"`
}

type AddToGroupRequest struct {
	ConversationId int64    `json:"conversationId"`
	Participants   []string `json:"participants"`
}

type LeaveGroupRequest struct {
	ConversationId int64 `json:"conversationId"`
}

type SetGroupNameRequest struct {
	ConversationId int64  `json:"conversationId"`
	Name           string `json:"name"`
}

type SetGroupPhotoRequest struct {
	ConversationId int64  `json:"conversationId"`
	PhotoId        string `json:"photoId"`
}

type SendMessageRequest struct {
	ReplyToMessageId *int64  `json:"replyTo"`
	Text             *string `json:"text,omitempty"`
	PhotoId          *string `json:"photoId,omitempty"`
}

type ForwardMessageRequest struct {
	ForwardToConversationId int64 `json:"forwardTo"`
	MessageId               int64 `json:"messageId"`
}

type ReactionRequest struct {
	Content string `json:"content"`
}
