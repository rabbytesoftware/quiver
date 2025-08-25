package database

import (
	"fmt"

	"github.com/rabbytesoftware/quiver/internal/logger"
)

// DatabaseType represents the type of database implementation
type DatabaseType string

const (
	DatabaseTypeJSON DatabaseType = "json"
)

// NewDatabase creates a new database instance based on the specified type
func NewDatabase(dbType DatabaseType, connectionString string, logger *logger.Logger) (Database, error) {
	switch dbType {
	case DatabaseTypeJSON:
		return NewJSONDatabase(connectionString, logger), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// NewDefaultDatabase creates a default JSON database instance
func NewDefaultDatabase(dbPath string, logger *logger.Logger) Database {
	return NewJSONDatabase(dbPath, logger)
} 