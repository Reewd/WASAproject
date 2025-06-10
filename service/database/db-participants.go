package database

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
