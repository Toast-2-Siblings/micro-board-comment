package database

import (
	"sync"

	"github.com/Toast-2-Siblings/micro-board-comment/config"
	"github.com/Toast-2-Siblings/micro-board-comment/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db_instance *gorm.DB
	once sync.Once
)

func InitDatabase() error {
	var err error
	cfg := config.GetConfig()
	
	if cfg.Mode == "production" {
		// production = postgresql
	} else { 
		// development = sqlite (in-memory)
		db_instance, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDB() *gorm.DB {
	once.Do(func() {
		if db_instance == nil {
			InitDatabase()
		}
	})
	return db_instance
}

func CloseDB() {
	if db_instance != nil {
		sqlDB, err := db_instance.DB()
		if err != nil {
			return
		}
		sqlDB.Close()
	}
}

func Migrate() error {
	return db_instance.AutoMigrate(
		&models.Comment{},
	)
}
