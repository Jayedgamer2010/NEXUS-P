//go:build sqlite

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDialector(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
