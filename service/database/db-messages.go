package database

import (
	"database/sql"
	"sort"

	"github.com/Reewd/WASAproject/service/database/helpers"
)

func (db *appdbimpl) InsertMessage(conversationId int64, userId int64, content *string, photoId *string, replyTo *int64, isForwarded bool) (int64, string, error) {
	stmt := `INSERT into messages (conversationId, senderId, content, photoId, replyTo, isForwarded) VALUES (?, ?, ?, ?, ?, ?) RETURNING id, timestamp`
	var timestamp string
	var messageId int64

	err := db.c.QueryRow(stmt, conversationId, userId, content, photoId, replyTo, isForwarded).Scan(&messageId, &timestamp)
	if err != nil {
		return 0, "", err
	}

	return messageId, timestamp, nil
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
func (db *appdbimpl) GetChat(conversationID int64) ([]MessageView, error) {
	const stmt = `
    SELECT 
        m.id                  AS messageId,
        m.content             AS messageText,
        m.conversationId,
        m.photoId             AS messagePhotoId,
        i.path                AS messagePhotoPath,
        m.isForwarded         AS isForwarded,
        m.replyTo,
        m.timestamp           AS messageTimestamp,
        u.id                  AS messageSenderId,
        u.username            AS messageSenderUsername,
        u.photoId             AS messageSenderPhotoId,
        ui.path               AS messageSenderPhotoPath,
        r.content             AS reactionContent,
        r.timestamp           AS reactionTimestamp,
        ru.id                 AS reactionSenderId,
        ru.username           AS reactionSenderUsername,
        ru.photoId            AS reactionSenderPhotoId,
        ri.path               AS reactionSenderPhotoPath
    FROM messages m
    LEFT JOIN users u  ON m.senderId    = u.id
    LEFT JOIN reactions r  ON m.id     = r.messageId
    LEFT JOIN users ru ON r.senderId   = ru.id
    LEFT JOIN images i ON m.photoId = i.uuid
    LEFT JOIN images ui on u.photoId = ui.uuid
    LEFT JOIN images ri on ru.photoId = ri.uuid
    WHERE m.conversationId = ?
    `

	// Separate query for message statuses
	const statusStmt = `
    SELECT messageId, status 
    FROM message_status 
    WHERE messageId IN (
        SELECT id FROM messages WHERE conversationId = ?
    )
    `

	statusRows, err := db.c.Query(statusStmt, conversationID)
	if err != nil {
		return nil, err
	}
	defer helpers.CloseRows(statusRows)

	statusMap := make(map[int64][]string)
	for statusRows.Next() {
		var messageID int64
		var status string
		if err := statusRows.Scan(&messageID, &status); err != nil {
			return nil, err
		}
		statusMap[messageID] = append(statusMap[messageID], status)
	}

	// Check for errors in status query
	if err := statusRows.Err(); err != nil {
		return nil, err
	}

	// Now query and process the messages
	rows, err := db.c.Query(stmt, conversationID)
	if err != nil {
		return nil, err
	}
	defer helpers.CloseRows(rows)

	// map for message aggregation
	msgMap := make(map[int64]*MessageView)

	for rows.Next() {
		var (
			messageID                 int64
			nsmessageText             sql.NullString
			convID                    int64
			isForwarded               bool
			nsMessagePhotoID          sql.NullString
			nsMessagePhotoPath        sql.NullString
			nrReplyTo                 sql.NullInt64
			messageTimestamp          string
			senderID                  int64
			senderUsername            string
			nsSenderPhotoID           sql.NullString
			nsSenderPhotoPath         sql.NullString
			nsReactionContent         sql.NullString
			nrReactionTimestamp       sql.NullString
			nsReactionSenderID        sql.NullInt64
			nsReactionSenderUsername  sql.NullString
			nsReactionSenderPhotoID   sql.NullString
			nsReactionSenderPhotoPath sql.NullString
		)

		if err := rows.Scan(
			&messageID,
			&nsmessageText,
			&convID,
			&nsMessagePhotoID,
			&nsMessagePhotoPath,
			&isForwarded,
			&nrReplyTo,
			&messageTimestamp,
			&senderID,
			&senderUsername,
			&nsSenderPhotoID,
			&nsSenderPhotoPath,
			&nsReactionContent,
			&nrReactionTimestamp,
			&nsReactionSenderID,
			&nsReactionSenderUsername,
			&nsReactionSenderPhotoID,
			&nsReactionSenderPhotoPath,
		); err != nil {
			return nil, err
		}

		// Build nullable pointers
		var photoID *string
		if nsMessagePhotoID.Valid {
			photoID = &nsMessagePhotoID.String
		}
		var photoPath *string
		if nsMessagePhotoPath.Valid {
			photoPath = &nsMessagePhotoPath.String
		}
		var replyTo *int64
		if nrReplyTo.Valid {
			replyTo = &nrReplyTo.Int64
		}

		var messageText *string
		if nsmessageText.Valid {
			messageText = &nsmessageText.String
		}

		// Create the MessageView if first time we see this message
		msg, ok := msgMap[messageID]
		if !ok {
			// build the sender's photo pointer
			var senderPhoto *Photo
			if nsSenderPhotoID.Valid && nsSenderPhotoPath.Valid {
				senderPhoto = &Photo{
					PhotoId: nsSenderPhotoID.String,
					Path:    nsSenderPhotoPath.String,
				}
			}

			var photo *Photo
			if photoID != nil && photoPath != nil {
				photo = &Photo{
					PhotoId: *photoID,
					Path:    *photoPath,
				}
			}

			msg = &MessageView{
				MessageId:      messageID,
				Text:           messageText,
				ConversationId: convID,
				Photo:          photo,
				ReplyTo:        replyTo,
				Timestamp:      messageTimestamp,
				SentBy: User{
					UserId:   senderID,
					Username: senderUsername,
					Photo:    senderPhoto,
				},
				Reactions:   []ReactionView{},
				IsForwarded: isForwarded,
			}
			msgMap[messageID] = msg
		}

		// Append reaction if there is one
		if nsReactionContent.Valid {
			// build reaction-sender photo pointer
			var rSenderPhoto *Photo
			if nsReactionSenderPhotoID.Valid && nsReactionSenderPhotoPath.Valid {
				rSenderPhoto = &Photo{
					PhotoId: nsReactionSenderPhotoID.String,
					Path:    nsReactionSenderPhotoPath.String,
				}
			}
			// build reaction timestamp pointer
			var rTs string
			if nrReactionTimestamp.Valid {
				rTs = nrReactionTimestamp.String
			}
			// build reaction-sender info
			var reactionSenderID int64
			var senderName string
			if nsReactionSenderID.Valid {
				reactionSenderID = nsReactionSenderID.Int64
			}
			if nsReactionSenderUsername.Valid {
				senderName = nsReactionSenderUsername.String
			}

			msg.Reactions = append(msg.Reactions, ReactionView{
				SentBy: User{
					UserId:   reactionSenderID,
					Username: senderName,
					Photo:    rSenderPhoto,
				},
				Content:   nsReactionContent.String,
				Timestamp: rTs,
			})
		}
	}

	// check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Determine each message's final status
	for id, statuses := range statusMap {
		m, exists := msgMap[id]
		if !exists {
			continue // Skip if the message doesn't exist
		}

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

	var out []MessageView
	for _, m := range msgMap {
		out = append(out, *m)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Timestamp < out[j].Timestamp
	})

	return out, nil
}

func (db *appdbimpl) ForwardMessage(messageIdToForward int64, conversationId int64, forwarderId int64) (messageId int64, timestamp string, content *string, photoId *string, err error) {
	stmt := `SELECT content, photoId FROM messages WHERE id = ?`
	var nsText sql.NullString
	var nsPhotoId sql.NullString
	err = db.c.QueryRow(stmt, messageIdToForward).Scan(&content, &nsPhotoId)
	if err != nil {
		return 0, "", nil, nil, err
	}

	if nsPhotoId.Valid {
		photoId = &nsPhotoId.String
	}

	if nsText.Valid {
		content = &nsText.String
	}

	forwardedMessageId, timestamp, err := db.InsertMessage(conversationId, forwarderId, content, photoId, nil, true)
	if err != nil {
		return 0, "", nil, nil, err
	}

	return forwardedMessageId, timestamp, content, photoId, nil
}

func (db *appdbimpl) GetConversationIdFromMessageId(messageId int64) (int64, error) {
	stmt := `SELECT conversationId FROM messages WHERE id = ?`
	var conversationId int64
	err := db.c.QueryRow(stmt, messageId).Scan(&conversationId)
	if err != nil {
		return 0, err
	}
	return conversationId, nil
}

func (db *appdbimpl) GetLastMessage(conversationId int64) (*MessageView, error) {

	empty, err := db.IsConversationEmpty(conversationId)
	if err != nil {
		return nil, err
	}
	if empty {
		return nil, nil
	}

	stmt := `SELECT m.id, m.content, m.photoId, i.path, m.replyTo, m.timestamp, u.id, u.username, u.photoId, ui.path, m.isForwarded
            FROM messages m
            LEFT JOIN images i ON m.photoId = i.uuid
            JOIN users u ON m.senderId = u.id
            LEFT JOIN images ui ON u.photoId = ui.uuid
            WHERE m.conversationId = ?
            ORDER BY m.timestamp DESC
            LIMIT 1`

	var msg MessageView
	var nsText sql.NullString
	var nsPhotoId sql.NullString
	var nsPhotoPath sql.NullString
	var nsReplyTo sql.NullInt64
	var nsSenderPhotoId sql.NullString
	var nsSenderPhotoPath sql.NullString
	var isForwarded bool

	err = db.c.QueryRow(stmt, conversationId).Scan(
		&msg.MessageId,
		&nsText,
		&nsPhotoId,
		&nsPhotoPath,
		&nsReplyTo,
		&msg.Timestamp,
		&msg.SentBy.UserId,
		&msg.SentBy.Username,
		&nsSenderPhotoId,
		&nsSenderPhotoPath,
		&isForwarded,
	)
	if err != nil {
		return nil, err
	}

	if nsText.Valid {
		msg.Text = &nsText.String
	}

	if nsPhotoId.Valid && nsPhotoPath.Valid {
		msg.Photo = &Photo{
			PhotoId: nsPhotoId.String,
			Path:    nsPhotoPath.String,
		}
	}

	if nsReplyTo.Valid {
		msg.ReplyTo = &nsReplyTo.Int64
	}

	if nsSenderPhotoId.Valid && nsSenderPhotoPath.Valid {
		msg.SentBy.Photo = &Photo{
			PhotoId: nsSenderPhotoId.String,
			Path:    nsSenderPhotoPath.String,
		}
	}
	msg.IsForwarded = isForwarded
	msg.Status = "sent" // Default status for last message preview

	return &msg, nil
}

func (db *appdbimpl) IsConversationEmpty(conversationId int64) (bool, error) {
	stmt := `SELECT EXISTS(SELECT 1 FROM messages WHERE conversationId = ?)`
	var count int
	err := db.c.QueryRow(stmt, conversationId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
