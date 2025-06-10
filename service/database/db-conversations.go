package database

func (db *appdbimpl) InsertConversation(name string, participants []string, isGroup bool, photo *string) (int64, error) {
	var conversationId int64
	if photo != nil {
		stmt := `INSERT INTO conversations (name, isGroup, photoId) VALUES (?, ?, ?)`
		result, err := db.c.Exec(stmt, name, isGroup, *photo)
		if err != nil {
			return 0, err
		}
		conversationId, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else {
		stmt := `INSERT INTO conversations (name, isGroup) VALUES (?, ?)`
		result, err := db.c.Exec(stmt, name, isGroup)
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

func (db *appdbimpl) GetParticipants(conversationId int64) ([]string, error) {
	stmt := `SELECT u.username FROM participants p
			 JOIN users u ON p.userId = u.id
			 WHERE p.conversationId = ?`
	rows, err := db.c.Query(stmt, conversationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []string
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return nil, err
		}
		participants = append(participants, username)
	}

	return participants, nil
}

func (db *appdbimpl) GetConversationsByUserId(userId int64) ([]Conversation, error) {
	stmt := `SELECT c.id, c.name, c.isGroup, c.photoId FROM conversations c
			 JOIN participants p ON c.id = p.conversationId
			 WHERE p.userId = ?`
	rows, err := db.c.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(&conv.ConversationId, &conv.Name, &conv.IsGroup, &conv.PhotoId)
		if err != nil {
			return nil, err
		}
		conv.Participants, err = db.GetParticipants(conv.ConversationId)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (db *appdbimpl) GetConversationById(conversationId int64) (*Conversation, error) {
	stmt := `SELECT c.id, c.name, c.isGroup, c.photoId FROM conversations c
			 WHERE c.id = ?`
	row := db.c.QueryRow(stmt, conversationId)

	var conv Conversation
	err := row.Scan(&conv.ConversationId, &conv.Name, &conv.IsGroup, &conv.PhotoId)
	if err != nil {
		return nil, err
	}

	conv.Participants, err = db.GetParticipants(conv.ConversationId)
	if err != nil {
		return nil, err
	}

	return &conv, nil
}

func (db *appdbimpl) ParticipantExists(conversationId int64, userId int64) (bool, error) {
	stmt := `SELECT EXISTS(SELECT 1 FROM participants WHERE conversationId = ? AND userId = ?)`
	var exists bool
	err := db.c.QueryRow(stmt, conversationId, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
