package api

type ConversationResponse struct {
	ConversationId int64    `json:"conversationId"`
	Title          string   `json:"title,omitempty"`
	Participants   []string `json:"participants"`
	IsGroup        bool     `json:"isGroup"`
	PhotoId        *string  `json:"photoId,omitempty"`
}
