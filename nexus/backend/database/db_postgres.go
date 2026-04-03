//go:build !sqlite

package database

import "gorm.io/driver/postgres"

func getDialector(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
