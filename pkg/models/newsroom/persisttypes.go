package newsroom

import (
	"github.com/joincivil/go-common-priv/pkg/models/article"
)

type Newsroom struct {
	ID 		 uint
	Name     string
	Address  string
	Articles []article.Article
}

type NewsroomPersister interface {
	CreateNewsroom(newsroom *Newsroom) error
	UpdateNewsroom(newsroom *Newsroom) error
	AddArticle(article *article.Article) error
	NewsroomByID(newsroomID uint) (*Newsroom, error)
	GetArticlesForNewsroom(newsroomID uint) ([]article.Article, error)
}
