package article

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
	"fmt"
	"encoding/json"
	"github.com/pkg/errors"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	// Could make this configurable later if needed
	maxOpenConns    = 50
	maxIdleConns    = 10
	connMaxLifetime = time.Second * 1800 // 30 mins
)

type ArticleGorm struct {
	gorm.Model
	BlockData        postgres.Jsonb
	ArticleMetadata  postgres.Jsonb
	NewsroomAddress  string
	IndexedTimestamp time.Time
	RawJSON          postgres.Jsonb
}

func(a *ArticleGorm) ConvertToArticle() (*Article, error) {
	article := &Article{}
	// if it fails it probably hasnt been added yet, do nothing
	blockdata := ethTypes.Receipt{}
	if err := json.Unmarshal(a.BlockData.RawMessage, &blockdata); err == nil {
		article.BlockData = blockdata
	}
		

	// if it fails it probably hasnt been added yet, do nothing
	metadata := ArticleMetadata{}
	if err := json.Unmarshal(a.ArticleMetadata.RawMessage, &metadata); err == nil {	
		article.ArticleMetadata = metadata
	}

	article.RawJSON = a.RawJSON.RawMessage
	article.ID = a.ID
	article.NewsroomAddress = a.NewsroomAddress
	article.IndexedTimestamp = a.IndexedTimestamp

	return article, nil
}

func(a *ArticleGorm) PopulateFromArticle(article *Article) error {
	metaJSON, metaerr := json.Marshal(article.ArticleMetadata)
	if metaerr != nil {
		return metaerr
	}
	a.ArticleMetadata = postgres.Jsonb{metaJSON}

	if article.BlockData.TxHash == (ethCommon.Hash{}) {
		blockJSON, blockerr := json.Marshal(article.BlockData)
		if blockerr != nil {
			return blockerr
		}
		a.BlockData = postgres.Jsonb{blockJSON}
	}

	a.NewsroomAddress = article.NewsroomAddress
	a.IndexedTimestamp = article.IndexedTimestamp
	a.RawJSON = postgres.Jsonb{article.RawJSON}
	a.ID = article.ID

	return nil
}


type ArticleGormPGPersister struct {
	DB *gorm.DB
	version string
}

func NewArticleGormPGPersister(host string, port int, user string, password string, dbname string)  (*ArticleGormPGPersister, error) {
	articleGormPGPersister := &ArticleGormPGPersister{}
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

func (p *ArticleGormPGPersister) ArticleByID(articleID uint) (*Article, error) {
	articleGorm := &ArticleGorm{}
	if err := p.DB.First(articleGorm, articleID).Error; err != nil {
		return nil, err
	}
	
	return articleGorm.ConvertToArticle()
}

func (p *ArticleGormPGPersister) CreateArticle(article *Article) error {
	metaJSON, err := json.Marshal(article.ArticleMetadata)
	if err != nil {
		return err
	}

	articleGorm := ArticleGorm{
		NewsroomAddress: article.NewsroomAddress,
		ArticleMetadata: postgres.Jsonb{metaJSON},
		RawJSON: postgres.Jsonb{article.RawJSON},
	}

	if err := p.DB.Create(&articleGorm).Error; err != nil {
		return err
	}

	article.ID = articleGorm.ID
	return nil
}

func (p *ArticleGormPGPersister) UpdateArticle(article *Article) error {
	articleGorm := ArticleGorm{}

	articleGorm.PopulateFromArticle(article)

	if err := p.DB.Save(&articleGorm).Error; err != nil {
		return err
	}
	return nil
}
