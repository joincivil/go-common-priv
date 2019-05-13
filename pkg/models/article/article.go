package article

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Article struct {
	gorm.Model
	BlockData        postgres.Jsonb
	ArticleMetadata  postgres.Jsonb
	NewsroomAddress  string
	IndexedTimestamp time.Time
	RawJson          postgres.Jsonb
}
