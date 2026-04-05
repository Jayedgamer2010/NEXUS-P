package database

import (
	"log"
	"os"
	"time"

	"nexus/backend/config"
	"nexus/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var db *gorm.DB
	var err error

	logLevel := logger.Silent
	if os.Getenv("APP_ENV") != "production" {
		logLevel = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	}

	if cfg.DBDriver == "postgres" && cfg.DatabaseURL != "" {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
		if err != nil {
			return err
		}
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.DBPath), gormConfig)
		if err != nil {
			return err
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	DB = db

	if err := runMigrations(); err != nil {
		log.Printf("Warning: some migrations failed: %v", err)
	}

	return nil
}

func runMigrations() error {
	// Disable foreign key checks for postgres during migration
	if DB.Dialector.Name() == "postgres" {
		DB.Exec("SET session_replication_role = replica")
	}

	migrations := []func() error{
		func() error { return DB.AutoMigrate(&models.User{}) },
		func() error { return DB.AutoMigrate(&models.Node{}) },
		func() error { return DB.AutoMigrate(&models.Egg{}) },
		func() error { return DB.AutoMigrate(&models.Allocation{}) },
		func() error { return DB.AutoMigrate(&models.Server{}) },
		func() error { return DB.AutoMigrate(&models.Ticket{}) },
		func() error { return DB.AutoMigrate(&models.CoinTransaction{}) },
	}

	var errs []error
	for _, m := range migrations {
		if err := m(); err != nil {
			errs = append(errs, err)
		}
	}

	if DB.Dialector.Name() == "postgres" {
		DB.Exec("SET session_replication_role = DEFAULT")
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
