package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"invites.cc/utils"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	SSLMode    string
}

// Create a new DWHConfig
func DWHConfig() *Config {
	return &Config{
		DBUser:     utils.GetEnv("DB_USER", "postgres"),
		DBPassword: utils.GetEnv("DB_PASSWORD", "postgres"),
		DBName:     utils.GetEnv("DB_NAME", "invites"),
		DBHost:     utils.GetEnv("DB_HOST", "localhost"),
		DBPort:     utils.GetEnv("DB_PORT", "5432"),
		SSLMode:    utils.GetEnv("DB_SSL_MODE", "disable"),
	}
}

// DBConnectionString returns the database connection string
func (c *Config) DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.SSLMode)
}

// ConnectDB attempts to connect to the database with retries
func ConnectDB(dsn string, maxRetries int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			// Test the connection
			sqlDB, err := db.DB()
			if err != nil {
				log.Printf("Failed to get database instance (attempt %d/%d): %v", i+1, maxRetries, err)
				continue
			}
			if err := sqlDB.Ping(); err != nil {
				log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
				continue
			}
			return db, nil
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
}
