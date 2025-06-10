package database

func (db *appdbimpl) InsertImage(uuid string, path string) error {
	stmt := `INSERT INTO images (uuid, path) VALUES (?, ?)`
	_, err := db.c.Exec(stmt, uuid, path)
	if err != nil {
		return err
	}
	return nil
}
