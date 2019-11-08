package newsroom

import (
	"time"

	carticle "github.com/joincivil/go-common/pkg/article"
)

// Newsroom is the representation of a newsroom used outside the persisters
type Newsroom struct {
	ID       uint
	Name     string
	Address  string
	Meta     *Meta
	Articles []carticle.Article
}

// Meta represents some arbitrary metadata for a newsroom. Use this for
// arbitrary flags and other configuration.
type Meta struct {
	// Index enables the indexing of newsroom content
	Index bool `json:"index"`
	// Claim enables the content claim creation of newsroom content
	Claim bool `json:"claim"`
}

// Persister is and interface for persisting newsrooms
type Persister interface {
	CreateNewsroom(newsroom *Newsroom) error
	UpdateNewsroom(newsroom *Newsroom) error
	AddArticle(newsroomID uint, article *carticle.Article) error
	Newsrooms() ([]*Newsroom, error)
	NewsroomByID(newsroomID uint) (*Newsroom, error)
	GetArticlesForNewsroom(newsroomID uint) ([]carticle.Article, error)
	GetArticlesForNewsroomIndexedSinceDate(newsroomID uint, date time.Time) ([]carticle.Article, error)
	GetLatestArticleForNewsroom(newsroomID uint) (*carticle.Article, error)
}
