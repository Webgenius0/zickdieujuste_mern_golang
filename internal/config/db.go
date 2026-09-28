package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(cfg *Config) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             2 * time.Second, // Slow SQL threshold (increased for remote DB)
			LogLevel:                  logger.Warn,     // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,            // Enable color
		},
	)

	dsn := cfg.Dsn
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         newLogger,
	})

	if err != nil {
		panic("failed to connect database")
	} else {
		println("Database connection successful")
	}

	// Setup Connection Pool to reduce remote DB latency (fixes SLOW SQL)
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db
}
