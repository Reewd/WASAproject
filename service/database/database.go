/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

type UserDatabase interface {
	Login(string) (int64, error)
}

// AppDatabase is the interface through which all DB operations are performed.
type AppDatabase interface {
	UserDatabase
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// NewAppDatabase reads initdb.sql, applies schema, and returns an AppDatabase.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	initdbBytes, err := os.ReadFile("./service/database/initdb.sql")
	if err != nil {
		return nil, fmt.Errorf("error reading initdb.sql: %w", err)
	}

	_, err = db.Exec(string(initdbBytes))
	if err != nil {
		return nil, fmt.Errorf("error executing initdb.sql: %w", err)
	}

	return &appdbimpl{c: db}, nil
}

// Ping checks if the DB connection is alive.
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
