package database

func (db *appdbimpl) NewUser(name string) (int, error) {
	db.c.Exec("INSERT INTO users (name) VALUES (?) RETURNING id", name)

}

