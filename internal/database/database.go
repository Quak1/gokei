package database

import (
	"context"
	"database/sql"
	"errors"
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

var initialCategoryID int32 = 0
var adminUserID int32 = 0

func InitialCategoryID() int32 {
	return initialCategoryID
}

func AdminUserID() int32 {
	return adminUserID
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

	err = runDBMigrations(dbConnection)
	if err != nil {
		dbConnection.Close()
		return nil, err
	}

	queries := store.NewQueriesWrapper(dbConnection)

	err = createAdminUser(queries)
	if err != nil {
		dbConnection.Close()
		return nil, err
	}

	err = createInitialCategory(queries)
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

func runDBMigrations(db *sql.DB) error {
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

func createAdminUser(queries *store.QueriesWrapper) error {
	user := store.CreateUserParams{
		Username: "admin",
		Name:     "admin",
	}

	adminUser, err := queries.GetUserByUsername(context.Background(), user.Username)
	if err == nil {
		adminUserID = adminUser.ID
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hash, err := utils.HashPassword("admin")
	if err != nil {
		return err
	}
	user.PasswordHash = hash

	newUser, err := queries.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}

	adminUserID = newUser.ID

	return nil
}

func createInitialCategory(queries *store.QueriesWrapper) error {
	category := store.CreateCategoryParams{
		Name:   "InitialBalance",
		Color:  "#123",
		Icon:   "B",
		UserID: adminUserID,
	}

	initialCategory, err := queries.GetCategoryByName(context.Background(), store.GetCategoryByNameParams{
		Name:   category.Name,
		UserID: category.UserID,
	})
	if err == nil {
		initialCategoryID = initialCategory.ID
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	newCategory, err := queries.CreateCategory(context.Background(), category)
	if err != nil {
		return err
	}

	initialCategoryID = newCategory.ID

	return nil
}
