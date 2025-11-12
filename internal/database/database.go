package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Quak1/gokei/internal/database/store"
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
	projectRoot, err := findProjectRoot()
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

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (go.mod not found)")
		}
		dir = parent
	}
}
