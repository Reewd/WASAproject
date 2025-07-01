package database

type User struct {
	UserId   int64
	Username string
	Photo    *Photo // optional, can be nil
}

type Conversation struct {
	ConversationId int64
	Name           string
	Participants   []User
	IsGroup        bool
	Photo          *Photo
}

type ReactionView struct {
	SentBy    User
	Content   string
	Timestamp string
}

type MessageView struct {
	MessageId      int64
	SentBy         User
	ConversationId int64
	Text           *string
	Timestamp      string
	Photo          *Photo
	Reactions      []ReactionView // aggregated reactions from rows sharing the same messageId
	ReplyTo        *int64
	Status         string // e.g., "sent", "delivered", "read"
	IsForwarded    bool   // indicates if the message was forwarded
}

type Photo struct {
	PhotoId string
	Path    string
}
