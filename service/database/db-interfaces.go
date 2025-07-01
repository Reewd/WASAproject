package database

import "database/sql"

type MessageDatabase interface {
	InsertMessage(conversationId int64, userId int64, content *string, photoId *string, replyTo *int64, isForwarded bool) (int64, string, error)
	RemoveMessage(messageId int64) error
	GetSenderId(messageId int64) (int64, error)
	GetChat(conversationID int64) ([]MessageView, error)
	GetConversationIdFromMessageId(messageId int64) (int64, error)
	ForwardMessage(messageIdToForward int64, conversationId int64, forwarderId int64) (messageId int64, timestamp string, content *string, photoId *string, err error)
	GetLastMessage(conversationId int64) (*MessageView, error)
	IsConversationEmpty(conversationId int64) (bool, error)
}

type ReactionDatabase interface {
	InsertReaction(messageId int64, userId int64, reaction string) error
	RemoveReaction(messageId int64, userId int64) error
	GetReactions(messageId int64) ([]ReactionView, error)
}

type ParticipantDatabase interface {
	InsertParticipants(conversationId int64, userId []int64) error
	RemoveParticipant(conversationId int64, userId int64) error
	GetParticipants(conversationId int64) ([]User, error)
	GetParticipantIds(conversationId int64) ([]int64, error)
}

type StatusDatabase interface {
	InsertSent(messageId int64, conversationId int64, recipientIds []int64) error
	InsertDelivered(recipientId int64) error
	InsertRead(conversationId int64, recipientId int64) error
}

type GroupDatabase interface {
	UpdateGroupName(conversationId int64, name string) error
	UpdateGroupPhoto(conversationId int64, photoId string) error
}

type ConversationDatabase interface {
	InsertConversation(name string, participants []string, isGroup bool, photo *string) (int64, error)
	GetConversationsByUserId(userId int64) ([]Conversation, error)
	GetConversationById(conversationId int64) (*Conversation, error)
	ParticipantExists(conversationId int64, userId int64) (bool, error)
	PrivateConversationExists(participants []string) (int64, error)
}

// All user related operations on the DB are handled by this interface.
type UserDatabase interface {
	Login(string) (*User, error)
	GetUserId(string) (int64, error)
	GetUsersIds([]string) ([]int64, error)
	GetUsername(int64) (string, error)
	InsertUser(string) (int64, error)
	UserExistsById(int64) (bool, error)
	UpdateUsername(string, int64) error
	UpdateUserPhoto(string, int64) error
	GetUsersByName([]string) ([]User, error)
	GetUser(int64) (*User, error)
	GetAllUsers() ([]User, error)
}

// All image related operations on the DB are handled by this interface.
type ImageDatabase interface {
	InsertImage(uuid string, path string) error
	GetImagePath(uuid string) (string, error)
}

// AppDatabase is the interface through which all DB operations are performed.
type AppDatabase interface {
	UserDatabase
	ImageDatabase
	ConversationDatabase
	ParticipantDatabase
	GroupDatabase
	MessageDatabase
	ReactionDatabase
	StatusDatabase
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}
