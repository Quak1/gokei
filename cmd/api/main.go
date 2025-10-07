package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	db *sql.DB
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4444, "Server port")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("GOKEI_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	requireFlag("dsn", cfg.db.dsn)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg.db.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := application{
		db: db,
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server on", "addr", srv.Addr)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
