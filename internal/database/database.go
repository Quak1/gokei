package database

import (
	"database/sql"
	"errors"

	"github.com/Quak1/gokei/internal/database/store"
)

type DB struct {
	Connection *sql.DB
	Queries    *store.Queries
}

var (
	ErrRecordNotFound = errors.New("record not found")
)

func OpenDB(dsn string) (*DB, error) {
	dbConnection, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = dbConnection.Ping()
	if err != nil {
		dbConnection.Close()
		return nil, err
	}

	queries := store.New(dbConnection)

	db := &DB{
		Connection: dbConnection,
		Queries:    queries,
	}

	return db, nil
}
