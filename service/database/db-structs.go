package database

type PublicUser struct {
	Username string
	PhotoId  *string
}

type Conversation struct {
	ConversationId int64
	Name           string
	Participants   []string
	IsGroup        bool
	PhotoId        *string
}

type ReactionView struct {
	SentBy    *PublicUser
	Content   *string
	Timestamp *int64
}

type MessageView struct {
	MessageId      int64
	SentBy         PublicUser
	ConversationId int64
	Content        string
	Timestamp      int64
	PhotoId        *string
	Reactions      []ReactionView // aggregated reactions from rows sharing the same messageId
	ReplyTo        *int64
	Status         string // e.g., "sent", "delivered", "read"
}
