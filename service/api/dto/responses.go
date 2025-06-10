package dto

type User struct {
	Username string `json:"name,omitempty"`
	UserId   int64  `json:"userId,omitempty"`
	PhotoId  string `json:"photoId,omitempty"`
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
