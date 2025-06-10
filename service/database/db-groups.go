package database

func (db *appdbimpl) UpdateGroupName(conversationId int64, name string) error {
	stmt := `UPDATE conversations SET name = ? WHERE id = ?`
	_, err := db.c.Exec(stmt, name, conversationId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) UpdateGroupPhoto(conversationId int64, photoId string) error {
	stmt := `UPDATE conversations SET photoId = ? WHERE id = ?`
	_, err := db.c.Exec(stmt, photoId, conversationId)
	if err != nil {
		return err
	}
	return nil
}
