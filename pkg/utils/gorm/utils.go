package gorm

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // need the driver
)

// NewGormPGConnection is a helper function to create a new Gorm conn pool given the
// Postgresql database parameters
func NewGormPGConnection(host string, port int, user string, password string,
	dbname string, maxOpenConns int, maxIdleConns int, connMaxLifetime time.Duration) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to gorm")
	}

	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetConnMaxLifetime(connMaxLifetime)

	return db, nil
}
