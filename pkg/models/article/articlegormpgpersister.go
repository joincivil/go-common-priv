package article

import (
	"encoding/json"
	"fmt"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"time"
)

const (
	// Could make this configurable later if needed
	maxOpenConns    = 50
	maxIdleConns    = 10
	connMaxLifetime = time.Second * 1800 // 30 mins
)

// Gorm is the article schema
type Gorm struct {
	gorm.Model
	BlockData        postgres.Jsonb
	ArticleMetadata  postgres.Jsonb
	NewsroomAddress  string
	IndexedTimestamp time.Time
	RawJSON          postgres.Jsonb
}

// TableName sets the name of the corresponding table in the db
func (Gorm) TableName() string {
	return "articles"
}

// ConvertToArticle returns the gorm struct as the public article struct
func (a *Gorm) ConvertToArticle() (*Article, error) {
	article := &Article{}
	// if it fails it probably hasnt been added yet, do nothing
	blockdata := ethTypes.Receipt{}
	if err := json.Unmarshal(a.BlockData.RawMessage, &blockdata); err == nil {
		article.BlockData = blockdata
	}

	// if it fails it probably hasnt been added yet, do nothing
	metadata := Metadata{}
	if err := json.Unmarshal(a.ArticleMetadata.RawMessage, &metadata); err == nil {
		article.ArticleMetadata = metadata
	}

	article.RawJSON = a.RawJSON.RawMessage
	article.ID = a.ID
	article.NewsroomAddress = a.NewsroomAddress
	article.IndexedTimestamp = a.IndexedTimestamp

	return article, nil
}

// PopulateFromArticle takes an article struct and maps its properties onto a gorm struct
func (a *Gorm) PopulateFromArticle(article *Article) error {
	metaJSON, metaerr := json.Marshal(article.ArticleMetadata)
	if metaerr != nil {
		return metaerr
	}
	a.ArticleMetadata = postgres.Jsonb{RawMessage: metaJSON}

	if article.BlockData.TxHash == (ethCommon.Hash{}) {
		blockJSON, blockerr := json.Marshal(article.BlockData)
		if blockerr != nil {
			return blockerr
		}
		a.BlockData = postgres.Jsonb{RawMessage: blockJSON}
	}

	a.NewsroomAddress = article.NewsroomAddress
	a.IndexedTimestamp = article.IndexedTimestamp
	a.RawJSON = postgres.Jsonb{RawMessage: article.RawJSON}
	a.ID = article.ID

	return nil
}

// GormPGPersister is a persister that uses gorm and postgres
type GormPGPersister struct {
	DB *gorm.DB
}

// NewGormPGPersister return a new persister
func NewGormPGPersister(host string, port int, user string, password string, dbname string) (*GormPGPersister, error) {
	articleGormPGPersister := &GormPGPersister{}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		return articleGormPGPersister, errors.Wrap(err, "error connecting to gorm")
	}
	articleGormPGPersister.DB = db

	articleGormPGPersister.DB.DB().SetMaxOpenConns(maxOpenConns)
	articleGormPGPersister.DB.DB().SetMaxIdleConns(maxIdleConns)
	articleGormPGPersister.DB.DB().SetConnMaxLifetime(connMaxLifetime)
	return articleGormPGPersister, nil
}

// ArticleByID finds an article by its ID
func (p *GormPGPersister) ArticleByID(articleID uint) (*Article, error) {
	articleGorm := &Gorm{}
	if err := p.DB.First(articleGorm, articleID).Error; err != nil {
		return nil, err
	}

	return articleGorm.ConvertToArticle()
}

// CreateArticle saves an article to the db
func (p *GormPGPersister) CreateArticle(article *Article) error {
	metaJSON, err := json.Marshal(article.ArticleMetadata)
	if err != nil {
		return err
	}

	articleGorm := Gorm{
		NewsroomAddress: article.NewsroomAddress,
		ArticleMetadata: postgres.Jsonb{RawMessage: metaJSON},
		RawJSON:         postgres.Jsonb{RawMessage: article.RawJSON},
	}

	if err := p.DB.Create(&articleGorm).Error; err != nil {
		return err
	}

	article.ID = articleGorm.ID
	return nil
}

// UpdateArticle saves updates to an article stuct
func (p *GormPGPersister) UpdateArticle(article *Article) error {
	articleGorm := Gorm{}

	if err := articleGorm.PopulateFromArticle(article); err != nil {
		return err
	}

	if err := p.DB.Save(&articleGorm).Error; err != nil {
		return err
	}

	return nil
}
