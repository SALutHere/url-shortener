package sqlite

import (
	"database/sql"
	"errors"
	"url-shortener/internal/lib/e"
	"url-shortener/internal/storage"

	"modernc.org/sqlite"
	sqlib "modernc.org/sqlite/lib"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (storage *Storage, err error) {
	const fn = "storage.sqlite.New"
	defer func() { err = e.WrapIfErr(fn, err) }()

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, err
	}

	if _, err = stmt.Exec(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (err error) {
	const fn = "storage.sqlite.SaveURL"
	defer func() { err = e.WrapIfErr(fn, err) }()

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(urlToSave, alias); err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlib.SQLITE_CONSTRAINT_UNIQUE {
				return storage.ErrURLExists
			}
		}

		return err
	}

	return nil
}

func (s *Storage) GetURL(alias string) (resURL string, err error) {
	const fn = "storage.sqlite.GetURL"
	defer func() { err = e.WrapIfErr(fn, err) }()

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", err
	}

	if err = stmt.QueryRow(alias).Scan(&resURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", err
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) (err error) {
	const fn = "storage.sqlite.DeleteURL"
	defer func() { err = e.WrapIfErr(fn, err) }()

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(alias); err != nil {
		return err
	}

	return nil
}
