package sqlite

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(addr string) (*gorm.DB, error) {
	cfg := gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	if os.Getenv("DEBUG") == "1" {
		cfg.Logger = logger.Default.LogMode(logger.Error)
	}
	db, err := gorm.Open(sqlite.Open(addr), &cfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}
