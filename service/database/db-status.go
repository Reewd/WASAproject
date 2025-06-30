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

func (db *appdbimpl) InsertDelivered(conversationId int64, recipientId int64) error {
	stmt := `
        UPDATE message_status 
        SET status = 'delivered' 
        WHERE recipientId = ? AND status = 'sent'`

	_, err := db.c.Exec(stmt, conversationId, recipientId)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) InsertRead(conversationId int64, recipientId int64) error {
	stmt := `
        UPDATE message_status 
        SET status = 'read' 
        WHERE conversationId = ? 
          AND recipientId = ? 
          AND status = 'delivered'`

	_, err := db.c.Exec(stmt, conversationId, recipientId)
	if err != nil {
		return err
	}

	return nil
}
