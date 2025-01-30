package app

import (
	"os"
	"testing"
)

func TestNewApp(t *testing.T) {
	// Get test database URLs from environment
	dbURI := os.Getenv("DATABASE_URL")
	if dbURI == "" {
		t.Skip("DATABASE_URL not set")
	}

	redisURI := os.Getenv("REDIS_URL")
	if redisURI == "" {
		t.Skip("REDIS_URL not set")
	}

	// Test app initialization
	app, err := NewApp(dbURI, redisURI)
	if err != nil {
		t.Fatalf("Failed to create new app: %v", err)
	}

	// Verify app is properly initialized
	if app == nil {
		t.Fatal("App should not be nil")
	}

	if app.Data == nil {
		t.Fatal("App.Data should not be nil")
	}

	if app.Data.Db == nil {
		t.Fatal("App.Data.Db should not be nil")
	}

	if app.Data.Redis == nil {
		t.Fatal("App.Data.Redis should not be nil")
	}

	// Test database connections
	if err := app.Data.Db.Ping(); err != nil {
		t.Fatalf("PostgreSQL connection failed: %v", err)
	} else {
		t.Log("PostgreSQL connection successful")
	}

	if err := app.Data.Redis.Ping(app.Data.Redis.Context()).Err(); err != nil {
		t.Fatalf("Redis connection failed: %v", err)
	} else {
		t.Log("Redis connection successful")
	}

	// Test cleanup
	if err := app.Close(); err != nil {
		t.Fatalf("Failed to close app: %v", err)
	}
}

func TestNewAppInvalidConnections(t *testing.T) {
	// Test with invalid PostgreSQL connection
	_, err := NewApp("postgres://invalid:5432/nonexistent", "redis://localhost:6379")
	if err == nil {
		t.Error("Expected error with invalid PostgreSQL connection")
	}

	// Test with invalid Redis connection
	_, err = NewApp("postgres://localhost:5432/testdb", "redis://invalid:6379")
	if err == nil {
		t.Error("Expected error with invalid Redis connection")
	}
}
