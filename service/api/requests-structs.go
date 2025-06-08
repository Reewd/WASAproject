package api

type UsernameRequest struct {
	Username string `json:"name"`
}

type PhotoRequest struct {
	PhotoId string `json:"photoid"`
}

type CreateConversationRequest struct {
	Title        string   `json:"title,omitempty"`
	Participants []string `json:"participants"`
	IsGroup      bool     `json:"isGroup"`
	PhotoId      *string  `json:"photoId,omitempty"`
}
