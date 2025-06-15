package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) InsertMessage(conversationId int64, userId int64, content string, photoId *string, replyTo *int64) error {
	stmt := `INSERT into messages (conversationId, userId, content, photoId, replyTo) VALUES (?, ?, ?, ?, ?)`
	_, err := db.c.Exec(stmt, conversationId, userId, content, photoId, replyTo)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetSenderId(messageId int64) (int64, error) {
	stmt := `SELECT senderId FROM messages WHERE id = ?`
	var senderId int64
	err := db.c.QueryRow(stmt, messageId).Scan(&senderId)
	if err != nil {
		return 0, err
	}
	return senderId, nil
}

func (db *appdbimpl) RemoveMessage(messageId int64) error {
	stmt := `DELETE FROM messages WHERE id = ?`
	_, err := db.c.Exec(stmt, messageId)
	if err != nil {
		return err
	}
	return nil
}

// GetMessageViews fetches all messages (with their reactions and status) for a given conversation.
//
//	TODO: Split this into smaller functions
func (db *appdbimpl) GetMessageViews(conversationID int64) ([]MessageView, error) {
	const stmt = `
    SELECT 
        m.id                  AS messageId,
        m.content             AS messageContent,
        m.conversationId,
        m.photoId             AS messagePhotoId,
        m.replyTo,
        m.timestamp           AS messageTimestamp,
        u.username            AS messageSenderUsername,
        u.photoId             AS messageSenderPhotoId,
        r.content             AS reactionContent,
        r.timestamp           AS reactionTimestamp,
        ru.username           AS reactionSenderUsername,
        ru.photoId            AS reactionSenderPhotoId,
        ms.status             AS messageStatus
    FROM messages m
    LEFT JOIN users u  ON m.userId    = u.id
    LEFT JOIN reactions r  ON m.id     = r.messageId
    LEFT JOIN users ru ON r.senderId   = ru.id
    LEFT JOIN message_status ms ON m.id = ms.messageId
    WHERE m.conversationId = ?
    ORDER BY m.timestamp ASC
    `

	rows, err := db.c.Query(stmt, conversationID)
	if err != nil {
		return nil, fmt.Errorf("querying messages: %w", err)
	}
	defer rows.Close()

	// map for message aggregation
	msgMap := make(map[int64]*MessageView)
	// map for collecting status strings per message
	statusMap := make(map[int64][]string)

	for rows.Next() {
		var (
			messageID                int64
			messageContent           string
			convID                   int64
			nsMessagePhotoID         sql.NullString
			nrReplyTo                sql.NullInt64
			messageTimestamp         int64
			senderUsername           string
			nsSenderPhotoID          sql.NullString
			nsReactionContent        sql.NullString
			nrReactionTimestamp      sql.NullInt64
			nsReactionSenderUsername sql.NullString
			nsReactionSenderPhotoID  sql.NullString
			MessageStatus            string
		)

		if err := rows.Scan(
			&messageID,
			&messageContent,
			&convID,
			&nsMessagePhotoID,
			&nrReplyTo,
			&messageTimestamp,
			&senderUsername,
			&nsSenderPhotoID,
			&nsReactionContent,
			&nrReactionTimestamp,
			&nsReactionSenderUsername,
			&nsReactionSenderPhotoID,
			&MessageStatus,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		// Build nullable pointers
		var photoID *string
		if nsMessagePhotoID.Valid {
			photoID = &nsMessagePhotoID.String
		}
		var replyTo *int64
		if nrReplyTo.Valid {
			replyTo = &nrReplyTo.Int64
		}

		statusMap[messageID] = append(statusMap[messageID], MessageStatus)

		// Create the MessageView if first time we see this message
		msg, ok := msgMap[messageID]
		if !ok {
			// build the sender’s photo pointer
			var senderPhoto *string
			if nsSenderPhotoID.Valid {
				senderPhoto = &nsSenderPhotoID.String
			}

			msg = &MessageView{
				MessageId:      messageID,
				Content:        messageContent,
				ConversationId: convID,
				PhotoId:        photoID,
				ReplyTo:        replyTo,
				Timestamp:      messageTimestamp,
				SentBy: PublicUser{
					Username: senderUsername,
					PhotoId:  senderPhoto,
				},
				Reactions: []ReactionView{},
			}
			msgMap[messageID] = msg
		}

		// Append reaction if there is one
		if nsReactionContent.Valid {
			// build reaction-sender photo pointer
			var rSenderPhoto *string
			if nsReactionSenderPhotoID.Valid {
				rSenderPhoto = &nsReactionSenderPhotoID.String
			}
			// build reaction timestamp pointer
			var rTs *int64
			if nrReactionTimestamp.Valid {
				rTs = &nrReactionTimestamp.Int64
			}
			// build reaction-sender username
			senderName := nsReactionSenderUsername.String

			msg.Reactions = append(msg.Reactions, ReactionView{
				SentBy: &PublicUser{
					Username: senderName,
					PhotoId:  rSenderPhoto,
				},
				Content:   &nsReactionContent.String,
				Timestamp: rTs,
			})
		}
	}

	// check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Determine each message’s final status
	for id, statuses := range statusMap {
		m := msgMap[id]

		allRead, allDelivered := true, true
		for _, s := range statuses {
			if s != "read" {
				allRead = false
			}
			if s != "delivered" && s != "read" {
				allDelivered = false
			}
			if !allRead && !allDelivered {
				break
			}
		}

		switch {
		case allRead:
			m.Status = "read"
		case allDelivered:
			m.Status = "delivered"
		default:
			m.Status = "sent"
		}
	}

	// Convert map→slice, then sort by timestamp
	var out []MessageView
	for _, m := range msgMap {
		out = append(out, *m)
	}

	return out, nil
}

func (db *appdbimpl) ForwardMessage(messageId int64, conversationId int64, forwarderId int64) error {
	stmt := `SELECT content, photoId FROM messages WHERE id = ?`
	var content string
	var photoId sql.NullString
	err := db.c.QueryRow(stmt, messageId).Scan(&content, &photoId)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO messages (conversationId, senderId, content, photoId) VALUES (?, ?, ?, ?)`
	_, err = db.c.Exec(stmt, conversationId, forwarderId, content, photoId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetConversationIdFromMessageId(messageId int64) (int64, error) {
	stmt := `SELECT conversationId FROM messages WHERE id = ?`
	var conversationId int64
	err := db.c.QueryRow(stmt, messageId).Scan(&conversationId)
	if err != nil {
		return 0, fmt.Errorf("getting conversation ID from message ID: %w", err)
	}
	return conversationId, nil
}
