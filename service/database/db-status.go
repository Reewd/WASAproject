package database

func (db *appdbimpl) InsertSent(messageId int64, conversationId int64, recipientIds []int64) error {
	status := "sent"
	stmt := `INSERT INTO message_status (messageId, conversationId, recipientId, status) VALUES (?, ?, ?, ?)`
	preparedStmt, err := db.c.Prepare(stmt)
	if err != nil {
		return err
	}
	defer preparedStmt.Close()

	for _, recipientId := range recipientIds {
		_, err := preparedStmt.Exec(messageId, conversationId, recipientId, status)
		if err != nil {
			return err
		}
	}
	return nil
}
