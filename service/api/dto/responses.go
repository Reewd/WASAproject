package dto

type User struct {
	Username string  `json:"name,omitempty"`
	UserId   int64   `json:"userId,omitempty"`
	PhotoId  *string `json:"photoId,omitempty"`
}

type PublicUser struct {
	Username string  `json:"username"`
	PhotoId  *string `json:"photoId,omitempty"`
}

type Photo struct {
	PhotoId string `json:"photoId"`
	Path    string `json:"path"`
}

type ConversationPreview struct {
	ConversationId int64        `json:"conversationId,omitempty"`
	Name           string       `json:"name,omitempty"`
	Participants   []PublicUser `json:"participants"`
	IsGroup        bool         `json:"isGroup"`
	PhotoId        *string      `json:"photoId,omitempty"`
	LastMessage    *SentMessage `json:"lastMessage,omitempty"` // optional, can be nil if no messages exist
}

type Chat struct {
	ConversationId int64         `json:"conversationId"`
	Name           string        `json:"name,omitempty"`
	Participants   []PublicUser  `json:"participants"`
	IsGroup        bool          `json:"isGroup"`
	PhotoId        *string       `json:"photoId,omitempty"`
	Messages       []SentMessage `json:"messages,omitempty"` // messages in the chat, can be empty if no messages exist
}

type Reaction struct {
	SentBy    PublicUser `json:"sentBy"`
	Content   string     `json:"content"`
	Timestamp string     `json:"timestamp"`
}

type SentMessage struct {
	MessageId        int64      `json:"messageId"`
	Content          *string    `json:"content"`
	ConversationId   int64      `json:"conversationId"` // ID of the conversation this message belongs to
	SentBy           PublicUser `json:"sentBy"`
	Timestamp        string     `json:"timestamp"`
	PhotoId          *string    `json:"photoId,omitempty"`
	Reactions        []Reaction `json:"reactions,omitempty"` // aggregated reactions from rows sharing the same messageId
	ReplyToMessageId *int64     `json:"replyTo,omitempty"`
	Status           string     `json:"status"` // e.g., "sent", "delivered", "read"
}
