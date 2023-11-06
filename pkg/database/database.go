package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nint8835/elf/pkg/config"
)

func Connect(config config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DatabasePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}

	err = db.AutoMigrate(&Guild{})
	if err != nil {
		return nil, fmt.Errorf("error migrating DB: %w", err)
	}

	return db, nil
}
