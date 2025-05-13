package store

import "database/sql"

type Store struct {
	db *sql.DB
}

// Constructor method patient store
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
