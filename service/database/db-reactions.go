package database

import (
	"database/sql"

	"github.com/Reewd/WASAproject/service/database/helpers"
)

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
	stmt := `
	SELECT 
		u.photoId,
		i.path,
		u.username, 
		r.content, 
		r.timestamp 
	FROM 
		reactions AS r 
	LEFT JOIN users AS u ON r.senderId = u.Id
	LEFT JOIN images AS i ON u.photoId = i.id
	WHERE r.messageId = ?
	`
	rows, err := db.c.Query(stmt, messageId)
	if err != nil {
		return nil, err
	}
	defer helpers.CloseRows(rows)

	var reactions []ReactionView
	for rows.Next() {
		var nsPhotoId sql.NullString
		var nsImagePath sql.NullString
		var reaction ReactionView

		if err := rows.Scan(&nsPhotoId, &nsImagePath, &reaction.SentBy.Username, &reaction.Content, &reaction.Timestamp); err != nil {
			return nil, err
		}
		if nsPhotoId.Valid {
			reaction.SentBy.Photo.PhotoId = nsPhotoId.String
		}
		if nsImagePath.Valid {
			reaction.SentBy.Photo.Path = nsImagePath.String
		}
		reactions = append(reactions, reaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}
