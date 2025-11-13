package testutils

import (
	"fmt"
	"log"
	"time"

	"github.com/Quak1/gokei/internal/database"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	// https://hub.docker.com/_/postgres
	POSTGRES_IMAGE   = "postgres"
	POSTGRES_VERSION = "17-alpine"

	POSTGRES_DB       = "testdb"
	POSTGRES_USER     = "postgres"
	POSTGRES_PASSWORD = "password"
)

var pool *dockertest.Pool

func NewTestDB() (*database.DB, func(), error) {
	var err error
	if pool == nil {
		// Uses a sensible default on windows (tcp/http) and linux/osx (socket).
		pool, err = dockertest.NewPool("")
		if err != nil {
			return nil, nil, fmt.Errorf("Could not construct pool: %v", err)
		}
	}

	err = pool.Client.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to connect to docker: %v", err)
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: POSTGRES_IMAGE,
			Tag:        POSTGRES_VERSION,
			Env: []string{
				"POSTGRES_USER=" + POSTGRES_USER,
				"POSTGRES_PASSWORD=" + POSTGRES_PASSWORD,
				"POSTGRES_DB=" + POSTGRES_DB,
				"listen_addresses = '*'",
			},
			Labels: map[string]string{"test_db": "1"},
		},
		func(config *docker.HostConfig) {
			// Set AutoRemove to true so that stopped container goes away by itself.
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not start resource: %v", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, hostAndPort, POSTGRES_DB)
	log.Println("Connecting to database on url: ", databaseUrl)

	var db *database.DB

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet.
	pool.MaxWait = 120 * time.Second
	if err := pool.Retry(func() error {
		db, err = database.OpenDB(databaseUrl)
		if err != nil {
			return err
		}
		return db.Connection.Ping()
	}); err != nil {
		return nil, nil, fmt.Errorf("Could not connect to docker db: %v", err)
	}

	cleanup := func() {
		if db != nil {
			db.Connection.Close()
		}
		if err := pool.Purge(resource); err != nil {
			log.Printf("Could not purge resource: %v", err)
		}
	}

	return db, cleanup, nil
}
