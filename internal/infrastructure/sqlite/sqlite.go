package sqlite

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/glebarez/sqlite"
)

func NewSQLiteDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	if err := db.AutoMigrate(
		&IpVisitStorage{},
		&IpHistoryStorage{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return db, nil
}