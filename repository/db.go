package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"loan-module/providers"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(dsn string) *Database {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	return &Database{DB: db}
}

// NewDatabaseWithConfig creates a new database connection with configuration
func NewDatabaseWithConfig(config *providers.Config) *Database {
	// Create the database connection
	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database connection: ", err)
	}

	// Set connection pool settings from config
	sqlDB.SetMaxIdleConns(config.DB.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.DB.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(config.DB.Timeout) * time.Minute)

	return &Database{DB: db}
}
