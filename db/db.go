package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type LinkStore struct {
	Db *sql.DB
}

func NewLinkStore(dbName string) (LinkStore, error) {
	Db, err := getConnection(dbName)
	if err != nil {
		return LinkStore{}, err
	}

	if err := createMigrations(Db); err != nil {
		return LinkStore{}, err
	}

	return LinkStore{
		Db,
	}, nil
}

func getConnection(dbName string) (*sql.DB, error) {
	var (
		err error
		db  *sql.DB
	)

	// Init SQLite3 database
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	log.Println("🚀 Connected Successfully to the Database")

	return db, nil
}

func createMigrations(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS visitors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		link_id INTEGER NOT NULL,
		agent VARCHAR(255) NOT NULL,
		created_at DATETIME default CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		full TEXT NOT NULL,
		short VARCHAR(64) NOT NULL,
		hits INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME default CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}
