package data

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
)

// Data encapsulates all database operations
type Data struct {
	db    *sql.DB
	redis *redis.Client
}

// NewData creates a new Data instance with database connections
func NewData(db *sql.DB, redis *redis.Client) *Data {
	return &Data{
		db:    db,
		redis: redis,
	}
}