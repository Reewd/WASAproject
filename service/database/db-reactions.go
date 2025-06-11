package database

func (db *appdbimpl) InsertReaction(messageId, senderId int64, content string) error {
	stmt := `INSERT INTO reactions (messageId, senderId, content) VALUES (?, ?, ?)`
	_, err := db.c.Exec(stmt, messageId, senderId, content)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveReaction(messageId, senderId int64) error {
	stmt := `DELETE FROM reactions WHERE messageId = ? AND senderId = ?`
	_, err := db.c.Exec(stmt, messageId, senderId)
	if err != nil {
		return err
	}
	return nil
}
