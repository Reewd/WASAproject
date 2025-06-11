package database

import "database/sql"

type MessageDatabase interface {
	InsertMessage(conversationId int64, userId int64, content string, photoId string) error
	InsertReply(conversationId int64, userId int64, content string, photoId string, replyTo int64) error
	RemoveMessage(messageId int64) error
	GetMessages(conversationId int64) ([]MessageView, error)
}

type ReactionDatabase interface {
	InsertReaction(messageId int64, userId int64, reaction string) error
	RemoveReaction(messageId int64, userId int64) error
}

type ParticipantDatabase interface {
	InsertParticipants(conversationId int64, userId []int64) error
	RemoveParticipant(conversationId int64, userId int64) error
}

type GroupDatabase interface {
	UpdateGroupName(conversationId int64, name string) error
	UpdateGroupPhoto(conversationId int64, photoId string) error
}

type ConversationDatabase interface {
	InsertConversation(name string, participants []string, isGroup bool, photo *string) (int64, error)
	GetParticipants(conversationId int64) ([]string, error)
	GetConversationsByUserId(userId int64) ([]Conversation, error)
	GetConversationById(conversationId int64) (*Conversation, error)
	ParticipantExists(conversationId int64, userId int64) (bool, error)
}

// All user related operations on the DB are handled by this interface.
type UserDatabase interface {
	Login(string) (int64, error)
	GetUserId(string) (int64, error)
	GetUsersIds([]string) ([]int64, error)
	GetUsername(int64) (string, error)
	InsertUser(string) (int64, error)
	UserIdExists(int64) (bool, error)
	UpdateUsername(string, int64) error
	UpdateUserPhoto(string, int64) error
}

// All image related operations on the DB are handled by this interface.
type ImageDatabase interface {
	InsertImage(uuid string, path string) error
}

// AppDatabase is the interface through which all DB operations are performed.
type AppDatabase interface {
	UserDatabase
	ImageDatabase
	ConversationDatabase
	ParticipantDatabase
	GroupDatabase
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}
