package helpers

import (
	"github.com/Reewd/WASAproject/service/api/dto"
	"github.com/Reewd/WASAproject/service/database"
)

func ConvertUser(user database.User) dto.User {
	return dto.User{
		UserId:   user.UserId,
		Username: user.Username,
		Photo:    ConvertPhoto(user.Photo),
	}
}

func ConvertUsers(users []database.User) []dto.User {
	dtoUsers := make([]dto.User, 0, len(users))
	for _, user := range users {
		dtoUsers = append(dtoUsers, ConvertUser(user))
	}
	return dtoUsers
}

func ConvertReactions(reactions []database.ReactionView) []dto.Reaction {
	convertedReactions := make([]dto.Reaction, 0, len(reactions))
	for _, reaction := range reactions {
		convertedReactions = append(convertedReactions, dto.Reaction{
			SentBy:    ConvertUser(reaction.SentBy),
			Content:   reaction.Content,
			Timestamp: reaction.Timestamp,
		})
	}
	return convertedReactions
}

func ConvertToSentMessages(messages []database.MessageView) []dto.SentMessage {
	sentMessages := make([]dto.SentMessage, 0, len(messages))
	for _, msg := range messages {
		sentMessages = append(sentMessages, dto.SentMessage{
			MessageId:        msg.MessageId,
			Text:             msg.Text,
			SentBy:           ConvertUser(msg.SentBy),
			Timestamp:        msg.Timestamp,
			Photo:            ConvertPhoto(msg.Photo),
			Reactions:        ConvertReactions(msg.Reactions),
			ReplyToMessageId: msg.ReplyTo,
			Status:           msg.Status,
			ConversationId:   msg.ConversationId,
			IsForwarded:      msg.IsForwarded,
		})
	}
	return sentMessages
}

func ConvertToSentMessage(msg database.MessageView) dto.SentMessage {
	return dto.SentMessage{
		MessageId:        msg.MessageId,
		Text:             msg.Text,
		SentBy:           ConvertUser(msg.SentBy),
		Timestamp:        msg.Timestamp,
		Photo:            ConvertPhoto(msg.Photo),
		Reactions:        ConvertReactions(msg.Reactions),
		ReplyToMessageId: msg.ReplyTo,
		Status:           msg.Status,
	}
}

func ConvertPhoto(photo *database.Photo) *dto.Photo {
	if photo == nil {
		return nil
	}
	return &dto.Photo{
		PhotoId: photo.PhotoId,
		Path:    photo.Path,
	}
}
