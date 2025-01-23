package api

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/rick/fly_crypto/internal/data"
)

// App holds application-wide dependencies
type App struct {
	data *data.Data
}

// NewApp creates a new App instance with database connections
func NewApp(dbURI, redisURI string) (*App, error) {
	// Initialize PostgreSQL connection
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, err
	}

	// Initialize Redis connection
	opts, err := redis.ParseURL(redisURI)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(opts)

	// Create new Data instance
	data := data.NewData(db, redisClient)

	return &App{
		data: data,
	}, nil
}

// Close closes all database connections
func (a *App) Close() error {
	if err := a.data.db.Close(); err != nil {
		return err
	}
	if err := a.data.redis.Close(); err != nil {
		return err
	}
	return nil
}