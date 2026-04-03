package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName      string
	AppPort      string
	AppSecret    string
	DBDriver     string
	DatabaseURL  string
	DBPath       string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPass       string
	JWTSecret    string
	JWTExpire    int
	WingsTokenID string
	WingsToken   string
}

func Load() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		AppName:      getEnv("APP_NAME", "NEXUS"),
		AppPort:      getEnv("APP_PORT", "3000"),
		AppSecret:    getEnv("APP_SECRET", "changeme"),
		DBDriver:     getEnv("DB_DRIVER", "sqlite"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		DBPath:       getEnv("DB_PATH", "./nexus.db"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBName:       getEnv("DB_NAME", "nexus"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPass:       getEnv("DB_PASS", ""),
		JWTSecret:    getEnv("JWT_SECRET", "changeme"),
		JWTExpire:    getEnvAsInt("JWT_EXPIRE_HOURS", 72),
		WingsTokenID: getEnv("WINGS_TOKEN_ID", ""),
		WingsToken:   getEnv("WINGS_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}
