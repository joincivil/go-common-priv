package newsroom

import (
	"github.com/joincivil/go-common-priv/pkg/models/article"
)

// Newsroom is the representation of a newsroom used outside the persisters
type Newsroom struct {
	ID       uint
	Name     string
	Address  string
	Articles []article.Article
}

// Persister is and interface for persisting newsrooms
type Persister interface {
	CreateNewsroom(newsroom *Newsroom) error
	UpdateNewsroom(newsroom *Newsroom) error
	AddArticle(article *article.Article) error
	NewsroomByID(newsroomID uint) (*Newsroom, error)
	GetArticlesForNewsroom(newsroomID uint) ([]article.Article, error)
}
