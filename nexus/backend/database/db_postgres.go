//go:build !sqlite

package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDialector(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
