package database

import "database/sql"

func (db *appdbimpl) Login(username string) (int64, error) {
	var id int64
	stmt := `SELECT id FROM users WHERE username = ?`
	err := db.c.QueryRow(stmt, username).Scan(&id)
	if err == sql.ErrNoRows {
		return db.Signup(username)
	}

	return id, err

}

func (db *appdbimpl) Signup(username string) (int64, error) {
	stmt := `INSERT INTO users (username) VALUES (?)`
	result, err := db.c.Exec(stmt, username)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *appdbimpl) UserIDExists(id int64) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
	err := db.c.QueryRow(stmt, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *appdbimpl) SetMyUsername(username string, id int64) error {
	stmt := `UPDATE users SET username = ? WHERE id = ?`
	_, err := db.c.Exec(stmt, username, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) SetMyPhoto(photoId string, id int64) error {
	stmt := `UPDATE users SET photoId = ? WHERE id = ?`
	_, err := db.c.Exec(stmt, photoId, id)
	if err != nil {
		return err
	}
	return nil
}
