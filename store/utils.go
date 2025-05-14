package store

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	db  *sql.DB
	rdb *redis.Client
}

// Constructor method patient store
func NewStore(db *sql.DB, rdb *redis.Client) *Store {
	return &Store{db: db, rdb: rdb}
}

type LoginResponse struct {
	UserID         uuid.UUID `json:"userid"`
	HashedPassword string    `json:"hashpassword"`
}
