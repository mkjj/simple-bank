package database

import (
	"fmt"
	"log"
	"simple_bank/config"
	"simple_bank/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DBHost, config.DBUser, config.DBPassword,
		config.DBName, config.DBPort, config.DBSSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.Account{},
		&models.Entry{},
		&models.Transfer{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	// Create composite index for transfers as defined in your schema
	DB.Exec(`
		CREATE INDEX IF NOT EXISTS transfers_from_to_account_id_idx 
		ON transfers (from_account_id, to_account_id);
	`)

	log.Println("Database migration completed")
	return nil
}
