package database

import "database/sql"

func (db *appdbimpl) InsertReaction(messageId, senderId int64, content string) error {
	stmt := `
    INSERT INTO reactions (messageId, senderId, content)
    VALUES (?, ?, ?)
    ON CONFLICT(messageId, senderId) DO UPDATE
      SET content   = excluded.content,
          timestamp = CURRENT_TIMESTAMP
    `
	_, err := db.c.Exec(stmt, messageId, senderId, content)
	return err
}

func (db *appdbimpl) RemoveReaction(messageId, senderId int64) error {
	stmt := `DELETE FROM reactions WHERE messageId = ? AND senderId = ?`
	_, err := db.c.Exec(stmt, messageId, senderId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetReactions(messageId int64) ([]ReactionView, error) {
	stmt := `SELECT u.photoId, r.senderId, r.content, r.timestamp FROM reactions AS r LEFT JOIN users AS u ON r.senderId = u.userId WHERE r.messageId = ?`
	rows, err := db.c.Query(stmt, messageId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []ReactionView
	for rows.Next() {
		var photoId sql.NullString
		var reaction ReactionView
		if err := rows.Scan(&photoId, &reaction.SentBy, &reaction.Content, &reaction.Timestamp); err != nil {
			return nil, err
		}
		if photoId.Valid {
			reaction.SentBy.PhotoId = &photoId.String
		}
		reactions = append(reactions, reaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}
