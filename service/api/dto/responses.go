package dto

type User struct {
	Username string `json:"name,omitempty"`
	UserId   int64  `json:"userId,omitempty"`
	PhotoId  string `json:"photoId,omitempty"`
}

type PublicUser struct {
	Username string  `json:"username"`
	PhotoId  *string `json:"photoId,omitempty"`
}

type Photo struct {
	PhotoId string `json:"photoId"`
	Path    string `json:"path"`
}

type Conversation struct {
	ConversationId int64    `json:"conversationId,omitempty"`
	Name           string   `json:"name,omitempty"`
	Participants   []string `json:"participants"`
	IsGroup        bool     `json:"isGroup"`
	PhotoId        *string  `json:"photoId,omitempty"`
}

type Reaction struct {
	SentBy    *PublicUser
	Content   string
	Timestamp *int64
}

type SentMessage struct {
	MessageId        int64       `json:"messageId"`
	Content          *string     `json:"content"`
	SentBy           PublicUser  `json:"sentBy"`
	Timestamp        string      `json:"timestamp"`
	PhotoId          *string     `json:"photoId,omitempty"`
	Reactions        []*Reaction `json:"reactions,omitempty"` // aggregated reactions from rows sharing the same messageId
	ReplyToMessageId *int64      `json:"replyTo,omitempty"`
	Status           string      `json:"status"` // e.g., "sent", "delivered", "read"
}

type ForwardedMessage struct {
	MessageId        int64       `json:"messageId"`
	Content          *string     `json:"content"`
	SentBy           PublicUser  `json:"sentBy"`
	Timestamp        string      `json:"timestamp"`
	PhotoId          *string     `json:"photoId,omitempty"`
	Reactions        []*Reaction `json:"reactions,omitempty"` // aggregated reactions from rows sharing the same messageId
	ReplyToMessageId *int64      `json:"replyTo,omitempty"`
	Status           string      `json:"status"` // e.g., "sent", "delivered", "read"
}
