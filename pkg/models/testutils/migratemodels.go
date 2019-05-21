package testutils

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joincivil/go-common-priv/pkg/models/article"
	"github.com/joincivil/go-common-priv/pkg/models/newsroom"
)

func MigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(&newsroom.NewsroomGorm{}, &article.ArticleGorm{}).Error
}
