package utils

import (
	"log"
	"os"

	"gorm.io/gorm"
)

// GetEnv returns the value of an environment variable or a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CheckDBConnection verifies the database connection and logs the result
func CheckDBConnection(db *gorm.DB, context string) error {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get database instance during %s: %v", context, err)
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		log.Printf("Failed to ping database during %s: %v", context, err)
		return err
	}

	log.Printf("Database connection verified during %s", context)
	return nil
}
