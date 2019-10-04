package article

import (
	carticle "github.com/joincivil/go-common/pkg/article"
)

// Persister an interface for persisting articles
type Persister interface {
	ArticleByID(articleID uint) (*carticle.Article, error)
	CreateArticle(article *carticle.Article) error
	UpdateArticle(article *carticle.Article) error
}
