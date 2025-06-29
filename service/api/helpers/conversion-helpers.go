package helpers

import (
	"github.com/Reewd/WASAproject/service/api/dto"
	"github.com/Reewd/WASAproject/service/database"
)

func ConvertPublicUser(user database.PublicUser) dto.PublicUser {
	return dto.PublicUser{
		Username: user.Username,
		Photo:    ConvertPhoto(user.Photo),
	}
}

func ConvertPublicUsers(users []database.PublicUser) []dto.PublicUser {
	publicUsers := make([]dto.PublicUser, 0, len(users))
	for _, user := range users {
		publicUsers = append(publicUsers, ConvertPublicUser(user))
	}
	return publicUsers
}

func ConvertReactions(reactions []database.ReactionView) []dto.Reaction {
	convertedReactions := make([]dto.Reaction, 0, len(reactions))
	for _, reaction := range reactions {
		convertedReactions = append(convertedReactions, dto.Reaction{
			SentBy:    ConvertPublicUser(reaction.SentBy),
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
			SentBy:           ConvertPublicUser(msg.SentBy),
			Timestamp:        msg.Timestamp,
			Photo:            ConvertPhoto(msg.Photo),
			Reactions:        ConvertReactions(msg.Reactions),
			ReplyToMessageId: msg.ReplyTo,
			Status:           msg.Status,
			ConversationId:   msg.ConversationId,
		})
	}
	return sentMessages
}

func ConvertToSentMessage(msg database.MessageView) dto.SentMessage {
	return dto.SentMessage{
		MessageId:        msg.MessageId,
		Text:             msg.Text,
		SentBy:           ConvertPublicUser(msg.SentBy),
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
