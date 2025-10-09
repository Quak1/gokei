package database

import (
	"database/sql"

	"github.com/Quak1/gokei/internal/database/queries"
)

type DB struct {
	Connection *sql.DB
	Queries    *queries.Queries
}

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

	queries := queries.New(dbConnection)

	db := &DB{
		Connection: dbConnection,
		Queries:    queries,
	}

	return db, nil
}
