package testutils

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // need postgres drivers
	"github.com/joincivil/go-common-priv/pkg/models/article"
	"github.com/joincivil/go-common-priv/pkg/models/newsroom"
)

// MigrateModels makes sure the db schema is up to date when the test runs
func MigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(&newsroom.Gorm{}, &article.Gorm{}).Error
}
