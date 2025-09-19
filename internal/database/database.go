package database

import (
	"database/sql"
	"log"

	"github.com/Quak1/gokei/internal/database/queries"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
	Queries *queries.Queries
}

func NewDBConnection(dbURL string) *DB {
	dbConnection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	return &DB{
		DB:      dbConnection,
		Queries: queries.New(dbConnection),
	}
}

func (db *DB) Close() error {
	return db.DB.Close()
}
