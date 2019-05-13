package newsroom

import (
	"github.com/jinzhu/gorm"
	"github.com/joincivil/go-common-priv/pkg/models/article"
)

type Newsroom struct {
	gorm.Model
	Name     string
	Address  string            `gorm:"unique;not null"`
	Articles []article.Article `gorm:"foreignkey:NewsroomAddress;association_foreignkey:Address"`
}
