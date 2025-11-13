package database

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/pkg/utils"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type DB struct {
	Connection *sql.DB
	Queries    store.QuerierTx
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

	err = RunDBMigrations(dbConnection)
	if err != nil {
		dbConnection.Close()
		return nil, err
	}

	queries := store.NewQueriesWrapper(dbConnection)

	err = queries.InsertInitialCategory(context.Background())
	if err != nil {
		dbConnection.Close()
		return nil, err
	}

	db := &DB{
		Connection: dbConnection,
		Queries:    queries,
	}

	return db, nil
}

func RunDBMigrations(db *sql.DB) error {
	projectRoot, err := utils.FindProjectRoot()
	if err != nil {
		return err
	}

	dir := filepath.Join(projectRoot, "sql", "migrations")
	fs := os.DirFS(dir)

	p, err := goose.NewProvider(goose.DialectPostgres, db, fs)
	if err != nil {
		return err
	}

	_, err = p.Up(context.Background())
	if err != nil {
		return err
	}

	return nil
}
