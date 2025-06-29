package database

import (
	"database/sql"

	"github.com/Reewd/WASAproject/service/database/helpers"
)

func (db *appdbimpl) Login(username string) (*User, error) {
	var id int64
	var nsPhotoId sql.NullString
	var nsImagePath sql.NullString
	stmt := `SELECT id, photoId, i.path FROM users LEFT JOIN images as i ON users.photoId = i.uuid WHERE username = ?`
	err := db.c.QueryRow(stmt, username).Scan(&id, &nsPhotoId, &nsImagePath)
	if err == sql.ErrNoRows {
		id, err := db.InsertUser(username)
		if err != nil {
			return nil, err
		}
		return &User{UserId: id, Username: username}, nil
	} else if err != nil {
		return nil, err
	}

	var photo *Photo
	if nsPhotoId.Valid && nsImagePath.Valid {
		photo = &Photo{PhotoId: nsPhotoId.String, Path: nsImagePath.String}
	}

	return &User{UserId: id, Username: username, Photo: photo}, nil
}

func (db *appdbimpl) InsertUser(username string) (int64, error) {
	stmt := `INSERT INTO users (username) VALUES (?)`
	result, err := db.c.Exec(stmt, username)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *appdbimpl) UserExistsById(id int64) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
	err := db.c.QueryRow(stmt, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *appdbimpl) GetUsername(id int64) (string, error) {
	var username string
	stmt := `SELECT username FROM users WHERE id = ?`
	err := db.c.QueryRow(stmt, id).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (db *appdbimpl) GetUserId(username string) (int64, error) {
	var id int64
	stmt := `SELECT id FROM users WHERE username = ?`
	err := db.c.QueryRow(stmt, username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *appdbimpl) GetUsersIds(usernames []string) ([]int64, error) {
	ids := make([]int64, 0, len(usernames))
	stmt := `SELECT id FROM users WHERE username = ?`
	for _, username := range usernames {
		var id int64
		err := db.c.QueryRow(stmt, username).Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (db *appdbimpl) GetUserPhoto(id int64) (string, error) {
	var photoId string
	stmt := `SELECT photoId FROM users WHERE id = ?`
	err := db.c.QueryRow(stmt, id).Scan(&photoId)
	if err != nil {
		return "", err
	}
	return photoId, nil
}

func (db *appdbimpl) UpdateUsername(username string, id int64) error {
	stmt := `UPDATE users SET username = ? WHERE id = ?`
	_, err := db.c.Exec(stmt, username, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) UpdateUserPhoto(photoId string, id int64) error {
	stmt := `UPDATE users SET photoId = ? WHERE id = ?`
	_, err := db.c.Exec(stmt, photoId, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetPublicUsersByName(usernames []string) ([]PublicUser, error) {
	var publicUsers []PublicUser
	stmt := `SELECT username, photoId FROM users WHERE username = ?`
	for _, username := range usernames {
		var user PublicUser
		var nsPhotoId sql.NullString
		err := db.c.QueryRow(stmt, username).Scan(&user.Username, &nsPhotoId)
		if err != nil {
			return nil, err
		}

		publicUsers = append(publicUsers, user)
	}
	return publicUsers, nil
}

func (db *appdbimpl) GetPublicUser(id int64) (*PublicUser, error) {
	var user PublicUser
	var nsPhotoId sql.NullString
	var nsImagePath sql.NullString
	stmt := `SELECT username, photoId, i.path FROM users WHERE id = ? LEFT JOIN images as i ON users.photoId = images.id`
	err := db.c.QueryRow(stmt, id).Scan(&user.Username, &nsPhotoId, &nsImagePath)
	if err != nil {
		return nil, err
	}

	if nsPhotoId.Valid && nsImagePath.Valid {
		user.Photo = &Photo{PhotoId: nsPhotoId.String, Path: nsImagePath.String}
	}

	return &user, nil
}

func (db *appdbimpl) GetAllPublicUsers() ([]PublicUser, error) {
	stmt := `SELECT username, photoId, i.path FROM users 
             LEFT JOIN images AS i ON users.photoId = i.uuid`
	rows, err := db.c.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer helpers.CloseRows(rows)

	var publicUsers []PublicUser
	for rows.Next() {
		var user PublicUser
		var nsPhotoId sql.NullString
		var nsImagePath sql.NullString

		err := rows.Scan(&user.Username, &nsPhotoId, &nsImagePath)
		if err != nil {
			return nil, err
		}

		if nsPhotoId.Valid && nsImagePath.Valid {
			user.Photo = &Photo{PhotoId: nsPhotoId.String, Path: nsImagePath.String}
		}

		publicUsers = append(publicUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return publicUsers, nil
}
