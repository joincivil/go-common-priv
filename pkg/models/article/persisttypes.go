package article

import (
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"time"
	"encoding/json"
)

type Contributor struct {
	Role string
	Name string
}

type Image struct {
	URL string
	Hash string
	H int
	W int
}

type ArticleMetadata struct {
	Title string
	RevisionContentHash string
	RevisionContentURL string
	CanonicalURL string
	Slug string
	Description string
	Contributors []Contributor
	Images []Image
	Tags []string
	PrimaryTag string
	RevisionDate time.Time
	OriginalPublishDate time.Time
	Opinion bool
	CivilSchemaVersion string
}

type Article struct {
	ID uint
	BlockData ethTypes.Receipt
	ArticleMetadata ArticleMetadata
	NewsroomAddress string
	IndexedTimestamp time.Time
	RawJSON json.RawMessage
}

type ArticlePersister interface {
	ArticleByID(articleID int) (*Article, error)
	CreateArticle(article *Article) error
	UpdateArticle(article *Article) error
}
