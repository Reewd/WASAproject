package database

type Conversation struct {
	ConversationId int64    `json:"conversationId,omitempty"`
	Name           string   `json:"name,omitempty"`
	Participants   []string `json:"participants"`
	IsGroup        bool     `json:"isGroup"`
	PhotoId        *string  `json:"photoId,omitempty"`
}
