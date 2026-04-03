package database

import (
	"fmt"
	"nexus/backend/config"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var dsn string

	switch cfg.DBDriver {
	case "postgres":
		dsn = cfg.DatabaseURL
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	default:
		dsn = cfg.DBPath
	}

	dialector := getDialector(dsn)

	// Production: disable logging for performance
	gormLogger := logger.Default.LogMode(logger.Silent)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:      gormLogger,
		PrepareStmt: true,
	})
	if err != nil {
		return err
	}

	// Configure connection pool for low RAM VPS
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Max 10 open connections, max 5 idle
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
