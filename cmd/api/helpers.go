package main

import (
	"database/sql"
	"fmt"
	"os"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func requireFlag[T comparable](name string, value T) {
	var zero T
	if value == zero {
		fmt.Printf("Error: -%s flag is required\n", name)
		os.Exit(1)
	}
}
