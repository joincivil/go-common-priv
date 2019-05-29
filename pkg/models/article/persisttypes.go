package article

import (
	"encoding/json"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"time"
)

// Contributor is someone who was part of creating an article
type Contributor struct {
	Role string
	Name string
}

// Image saves information about images in the article including a potential content hash that can be used for verifications
type Image struct {
	URL  string
	Hash string
	H    int
	W    int
}

// ArticleMetadata is the information about an article that is indexeds
type ArticleMetadata struct {
	Title               string
	RevisionContentHash string
	RevisionContentURL  string
	CanonicalURL        string
	Slug                string
	Description         string
	Contributors        []Contributor
	Images              []Image
	Tags                []string
	PrimaryTag          string
	RevisionDate        time.Time
	OriginalPublishDate time.Time
	Opinion             bool
	CivilSchemaVersion  string
}

// Article is the top level representation of an article
type Article struct {
	ID               uint
	BlockData        ethTypes.Receipt
	ArticleMetadata  ArticleMetadata
	NewsroomAddress  string
	IndexedTimestamp time.Time
	RawJSON          json.RawMessage
}

// Persister an interface for persisting articles
type Persister interface {
	ArticleByID(articleID int) (*Article, error)
	CreateArticle(article *Article) error
	UpdateArticle(article *Article) error
}
