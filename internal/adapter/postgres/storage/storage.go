package storage

import (
	"database/sql"
)

type Storage struct {
	DB *sql.DB
}

func New(dbPath string) (*Storage, error) {
	db, err := openPostgres(dbPath)
	if err != nil {
		return nil, err
	}

	st := &Storage{
		DB: db,
	}

	if err := st.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return st, nil
}

func (s *Storage) Close() error {
	return s.DB.Close()
}