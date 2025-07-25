package dto

type LoginRequest struct {
	Username string `json:"username"`
}

type SetUsernameRequest struct {
	Username string `json:"username"`
}

type PhotoRequest struct {
	Photo *Photo `json:"photo,omitempty"`
}

type CreateConversationRequest struct {
	Name         string   `json:"name,omitempty"`
	Participants []string `json:"participants"`
	IsGroup      bool     `json:"isGroup"`
	Photo        *Photo   `json:"photo,omitempty"`
}

type AddToGroupRequest struct {
	ConversationId int64    `json:"conversationId"`
	Participants   []string `json:"participants"`
}

type LeaveGroupRequest struct {
	ConversationId int64 `json:"conversationId"`
}

type SetGroupNameRequest struct {
	Name string `json:"name"`
}

type SetGroupPhotoRequest struct {
	Photo *Photo `json:"photo,omitempty"`
}

type SendMessageRequest struct {
	ReplyToMessageId *int64  `json:"replyTo"`
	Text             *string `json:"text,omitempty"`
	Photo            *Photo  `json:"photo,omitempty"`
}

type ForwardMessageRequest struct {
	ForwardToConversationId int64 `json:"forwardTo"`
	MessageId               int64 `json:"messageId"`
}

type ReactionRequest struct {
	Content string `json:"content"`
}
