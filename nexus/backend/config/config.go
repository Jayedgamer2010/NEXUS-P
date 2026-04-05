package config

import (
	"os"
)

type Config struct {
	AppName       string
	AppPort       string
	DBDriver      string
	DatabaseURL   string
	DBPath        string
	JWTSecret     string
	JWTExpire     int
	AdminEmail    string
	AdminPassword string
	AdminUsername string
}

func Load() *Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	name := os.Getenv("APP_NAME")
	if name == "" {
		name = "NEXUS"
	}
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite"
	}
	dbURL := os.Getenv("DATABASE_URL")
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./nexus.db"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "change-this-to-random-32-char-string"
	}
	expire := 72
	return &Config{
		AppName:       name,
		AppPort:       port,
		DBDriver:      driver,
		DatabaseURL:   dbURL,
		DBPath:        dbPath,
		JWTSecret:     jwtSecret,
		JWTExpire:     expire,
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		AdminUsername: os.Getenv("ADMIN_USERNAME"),
	}
}
