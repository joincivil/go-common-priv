package newsroom

import (
	"time"

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
	AddArticle(newsroomID uint, article *article.Article) error
	Newsrooms() ([]*Newsroom, error)
	NewsroomByID(newsroomID uint) (*Newsroom, error)
	GetArticlesForNewsroom(newsroomID uint) ([]article.Article, error)
	GetArticlesForNewsroomIndexedSinceDate(newsroomID uint, date time.Time) ([]article.Article, error)
	GetLatestArticleForNewsroom(newsroomID uint) (*article.Article, error)
}
