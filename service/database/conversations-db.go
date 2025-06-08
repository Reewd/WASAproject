package database

func (db *appdbimpl) InsertConversation(title string, participants []string, isGroup bool, photo *string) (int64, error) {
	var conversationId int64
	if photo != nil {
		stmt := `INSERT INTO conversations (title, isGroup, photoId) VALUES (?, ?, ?)`
		result, err := db.c.Exec(stmt, title, isGroup, *photo)
		if err != nil {
			return 0, err
		}
		conversationId, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else {
		stmt := `INSERT INTO conversations (title, isGroup) VALUES (?, ?)`
		result, err := db.c.Exec(stmt, title, isGroup)
		if err != nil {
			return 0, err
		}
		conversationId, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	}

	err := db.InsertParticipantsFromUsername(conversationId, participants)
	if err != nil {
		return 0, err
	}
	return conversationId, nil
}

//TODO: Figure out how to handle the case where a participant does not exist in the users table.

func (db *appdbimpl) InsertParticipantsFromUsername(conversationId int64, participants []string) error {
	stmt := `INSERT INTO participants (conversationId, userId) VALUES (?, ?)`
	for _, participant := range participants {
		var userId int64
		err := db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, participant).Scan(&userId)
		if err != nil {
			return err
		}
		_, err = db.c.Exec(stmt, conversationId, userId)
		if err != nil {
			return err
		}
	}
	return nil
}
