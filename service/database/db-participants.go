package database

import (
	"database/sql"

	"github.com/Reewd/WASAproject/service/database/helpers"
)

func (db *appdbimpl) InsertParticipants(conversationId int64, userId []int64) error {
	stmt := `INSERT INTO participants (conversation_id, user_id) VALUES (?, ?)`
	for _, id := range userId {
		_, err := db.c.Exec(stmt, conversationId, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *appdbimpl) RemoveParticipant(conversationId int64, userId int64) error {
	stmt := `DELETE FROM participants WHERE conversation_id = ? AND user_id = ?`
	_, err := db.c.Exec(stmt, conversationId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetParticipants(conversationId int64) ([]PublicUser, error) {
	stmt := `SELECT u.username, u.photoId FROM participants p
		 JOIN users u ON p.userId = u.id
		 WHERE p.conversationId = ?`
	rows, err := db.c.Query(stmt, conversationId)
	if err != nil {
		return nil, err
	}
	defer helpers.CloseRows(rows)

	var participants []PublicUser
	for rows.Next() {
		var participant PublicUser
		var nsPhotoId sql.NullString
		err := rows.Scan(&participant.Username, &nsPhotoId)
		if err != nil {
			return nil, err
		}
		if nsPhotoId.Valid {
			participant.PhotoId = &nsPhotoId.String
		}
		participants = append(participants, participant)
	}

	// Check rows.Err after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return participants, nil
}

func (db *appdbimpl) GetParticipantIds(conversationId int64) ([]int64, error) {
	stmt := `SELECT userId FROM participants WHERE conversationId = ?`
	rows, err := db.c.Query(stmt, conversationId)
	if err != nil {
		return nil, err
	}
	defer helpers.CloseRows(rows)

	var participantIds []int64
	for rows.Next() {
		var userId int64
		err := rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		participantIds = append(participantIds, userId)
	}

	// Check rows.Err after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return participantIds, nil
}
