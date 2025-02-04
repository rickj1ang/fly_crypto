package app

import (
	"database/sql"
	"sync"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/rickj1ang/fly_crypto/internal/data"
)

// App holds application-wide dependencies
type App struct {
	Data         *data.Data
	SupportCoins []string
	CoinsPrices  *sync.Map
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

	var m sync.Map

	return &App{
		Data:         data,
		SupportCoins: []string{"BTC", "ETH", "SOL"},
		CoinsPrices:  &m,
	}, nil
}

// Close closes all database connections
func (a *App) Close() error {
	if err := a.Data.Db.Close(); err != nil {
		return err
	}
	if err := a.Data.Redis.Close(); err != nil {
		return err
	}
	return nil
}
