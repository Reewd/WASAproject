package database

func (db *appdbimpl) InsertImage(uuid string, path string) error {
	stmt := `INSERT INTO images (uuid, path) VALUES (?, ?)`
	_, err := db.c.Exec(stmt, uuid, path)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetImagePath(uuid string) (string, error) {
	var path string
	stmt := `SELECT path FROM images WHERE uuid = ?`
	err := db.c.QueryRow(stmt, uuid).Scan(&path)
	if err != nil {
		return "", err
	}
	return path, nil
}
