package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id         INTEGER  PRIMARY KEY AUTOINCREMENT,
			title      TEXT     NOT NULL,
			author     TEXT     NOT NULL,
			isbn       TEXT     DEFAULT '',
			genre      TEXT     DEFAULT '',
			cover_url  TEXT     DEFAULT '',
			date_read  TEXT     DEFAULT '',
			rating     INTEGER  DEFAULT 0 CHECK(rating BETWEEN 0 AND 5),
			review     TEXT     DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}
