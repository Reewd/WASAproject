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
